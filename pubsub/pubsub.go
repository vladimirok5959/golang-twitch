// Package implements Twitch API PubSub and automatically take care of API
// limits. Also it will handle automatically reconnections, ping/pong and
// maintenance requests.
package pubsub

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"
)

// Default Twitch server API credentials.
//
// https://dev.twitch.tv/docs/pubsub/#connection-management
const TwitchApiScheme = "wss"
const TwitchApiHost = "pubsub-edge.twitch.tv"
const TwitchApiPath = ""

const TwitchApiMaxTopics = 50

// PubSub is represent of API client.
type PubSub struct {
	sync.RWMutex

	URL         url.URL
	Connections map[int64]*Connection

	// Events
	eventOnConnect    func(*Connection)
	eventOnDisconnect func(*Connection)
	eventOnError      func(*Connection, error)
	eventOnInfo       func(*Connection, string)
	eventOnMessage    func(*Connection, *Answer)
	eventOnPing       func(*Connection, time.Time)
	eventOnPong       func(*Connection, time.Time, time.Time)
}

// New create and returns new API client.
func New() *PubSub {
	return NewWithURL(url.URL{
		Scheme: TwitchApiScheme,
		Host:   TwitchApiHost,
		Path:   TwitchApiPath,
	})
}

// NewWithURL create and returns new API client with custom API server URL.
// It can be useful for testing.
func NewWithURL(url url.URL) *PubSub {
	p := PubSub{
		URL:         url,
		Connections: map[int64]*Connection{},
	}
	return &p
}

// -----------------------------------------------------------------------------

func (p *PubSub) newConnection() *Connection {
	c := NewConnection(p.URL)
	c.OnConnect(p.eventOnConnect)
	c.OnDisconnect(p.eventOnDisconnect)
	c.OnError(p.eventOnError)
	c.OnInfo(p.eventOnInfo)
	c.OnMessage(p.eventOnMessage)
	c.OnPing(p.eventOnPing)
	c.OnPong(p.eventOnPong)
	return c
}

// -----------------------------------------------------------------------------

// Listen is adding topics for listening. It take care of API limits.
// New TCP connection will be created for every 50 topics.
//
// https://dev.twitch.tv/docs/pubsub/#connection-management
func (p *PubSub) Listen(ctx context.Context, topic string, params ...interface{}) {
	p.Lock()
	defer p.Unlock()

	// Create and add first connection
	if len(p.Connections) <= 0 {
		c := p.newConnection()
		p.Connections[c.ID] = c
	}

	t := p.Topic(topic, params...)

	// Check topic in connection
	// Don't continue if already present
	for _, c := range p.Connections {
		select {
		case <-ctx.Done():
			return
		default:
		}

		if c.HasTopic(t) {
			return
		}
	}

	// Add topic to first not busy connection
	for _, c := range p.Connections {
		select {
		case <-ctx.Done():
			return
		default:
		}

		if c.TopicsCount() < TwitchApiMaxTopics {
			c.AddTopic(t)
			return
		}
	}

	// Create new one and add
	c := p.newConnection()
	p.Connections[c.ID] = c
	c.AddTopic(t)
}

// Unlisten is remove topics from listening. It take care of API limits too.
// Connection count will automatically decrease of needs.
//
// https://dev.twitch.tv/docs/pubsub/#connection-management
func (p *PubSub) Unlisten(ctx context.Context, topic string, params ...interface{}) {
	p.Lock()
	defer p.Unlock()

	t := p.Topic(topic, params...)

	// Search and unlisten
	for _, c := range p.Connections {
		select {
		case <-ctx.Done():
			return
		default:
		}

		if c.HasTopic(t) {
			c.RemoveTopic(t)

			// Must not contain duplicates
			// So just remove first and break
			break
		}
	}

	// Remove empty connections
	for i, c := range p.Connections {
		select {
		case <-ctx.Done():
			return
		default:
		}

		if c.TopicsCount() <= 0 {
			_ = c.Close()
			delete(p.Connections, i)

			// Can't be more than one connection without topics
			// So just close, remove first and break
			break
		}
	}
}

// Topics returns all current listen topics.
func (p *PubSub) Topics() []string {
	p.Lock()
	defer p.Unlock()

	topics := []string{}
	for _, c := range p.Connections {
		for topic := range c.topics {
			topics = append(topics, topic)
		}
	}

	return topics
}

// HasTopic returns true if topic present.
func (p *PubSub) HasTopic(topic string, params ...interface{}) bool {
	p.Lock()
	defer p.Unlock()

	t := p.Topic(topic, params...)

	for _, c := range p.Connections {
		if c.HasTopic(t) {
			return true
		}
	}

	return false
}

// TopicsCount return count of topics.
func (p *PubSub) TopicsCount() int {
	p.Lock()
	defer p.Unlock()

	count := 0

	for _, c := range p.Connections {
		count += c.TopicsCount()
	}

	return count
}

// Topic generate correct topic for API.
// Params can be as number or string.
//
// https://dev.twitch.tv/docs/pubsub/#topics
func (p *PubSub) Topic(topic string, params ...interface{}) string {
	if len(params) <= 0 {
		return topic
	}

	var list []string
	for _, param := range params {
		list = append(list, fmt.Sprint(param))
	}

	return fmt.Sprintf("%s.%s", topic, strings.Join(list, "."))
}

// Close is close all connections.
// Usually need to call at the end of app life.
func (p *PubSub) Close() {
	p.Lock()
	defer p.Unlock()

	for i, c := range p.Connections {
		_ = c.Close()
		delete(p.Connections, i)
	}
}

// -----------------------------------------------------------------------------

// OnConnect is bind func to event.
// Will fire for every connection.
func (c *PubSub) OnConnect(fn func(*Connection)) {
	c.eventOnConnect = fn
}

// OnDisconnect is bind func to event.
// Will fire for every connection.
func (c *PubSub) OnDisconnect(fn func(*Connection)) {
	c.eventOnDisconnect = fn
}

// OnError is bind func to event.
// Will fire for every connection.
func (c *PubSub) OnError(fn func(*Connection, error)) {
	c.eventOnError = fn
}

// OnInfo is bind func to event.
// Will fire for every connection.
func (c *PubSub) OnInfo(fn func(*Connection, string)) {
	c.eventOnInfo = fn
}

// OnMessage is bind func to event.
// Will fire for every connection.
func (c *PubSub) OnMessage(fn func(*Connection, *Answer)) {
	c.eventOnMessage = fn
}

// OnPing is bind func to event.
// Will fire for every connection.
func (c *PubSub) OnPing(fn func(*Connection, time.Time)) {
	c.eventOnPing = fn
}

// OnPong is bind func to event.
// Will fire for every connection.
func (c *PubSub) OnPong(fn func(*Connection, time.Time, time.Time)) {
	c.eventOnPong = fn
}
