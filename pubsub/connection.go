package pubsub

import (
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// At least once every 5 minutes by docs but better keep this at 30 seconds
//
// https://dev.twitch.tv/docs/pubsub/#connection-management
const TwitchApiPingEach = 15 * time.Second    // 30 seconds
const TwitchApiPingTimeout = 10 * time.Second // 10 seconds

var nextConnectionID int64 = 0

// Connection is represent of one connection.
type Connection struct {
	sync.RWMutex

	done   chan struct{}
	topics map[string]*struct{}
	active bool
	url    url.URL

	ping_start  time.Time
	ping_sended bool

	ID         int64
	Connection *websocket.Conn

	// Events
	eventOnConnect    func(*Connection)
	eventOnDisconnect func(*Connection)
	eventOnError      func(*Connection, error)
	eventOnInfo       func(*Connection, string)
	eventOnMessage    func(*Connection, []byte)
	eventOnPing       func(*Connection, time.Time)
	eventOnPong       func(*Connection, time.Time, time.Time)
}

// NewConnection create new connection.
// Returns pointer to connection.
func NewConnection(url url.URL) *Connection {
	c := &Connection{
		done:   make(chan struct{}),
		topics: map[string]*struct{}{},
		active: false,
		url:    url,

		ping_start:  time.Now(),
		ping_sended: false,

		ID:         nextConnectionID,
		Connection: nil,
	}

	go_reconnector(c)
	go_reader(c)
	go_pinger(c)

	nextConnectionID++

	return c
}

// -----------------------------------------------------------------------------

func (c *Connection) onConnect() {
	if c.eventOnConnect != nil {
		c.eventOnConnect(c)
	}
}

func (c *Connection) onDisconnect() {
	if c.eventOnDisconnect != nil {
		c.eventOnDisconnect(c)
	}
}

func (c *Connection) onError(err error) {
	if c.eventOnError != nil {
		c.eventOnError(c, err)
	}
}

func (c *Connection) onInfo(str string) {
	if c.eventOnInfo != nil {
		c.eventOnInfo(c, str)
	}
}

func (c *Connection) onMessage(msg []byte) {
	if c.eventOnMessage != nil {
		c.eventOnMessage(c, msg)
	}
}

func (c *Connection) onPing(start time.Time) {
	if c.eventOnPing != nil {
		c.eventOnPing(c, start)
	}
}

func (c *Connection) onPong(start, end time.Time) {
	if c.eventOnPong != nil {
		c.eventOnPong(c, start, end)
	}
}

// -----------------------------------------------------------------------------

func (c *Connection) listenTopis() {
	topics := []string{}
	for topic := range c.topics {
		topics = append(topics, topic)
	}
	msg := `{"type":"LISTEN","nonce":"","data":{"topics":["` + strings.Join(topics, `","`) + `"]}}`
	if err := c.Connection.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		c.onError(err)
		c.active = false
		c.onDisconnect()
	}
}

// -----------------------------------------------------------------------------

// AddTopic is adding topics for listening.
func (c *Connection) AddTopic(topic string) {
	if _, ok := c.topics[topic]; ok {
		return
	}

	if len(c.topics) >= TwitchApiMaxTopics {
		return
	}

	c.Lock()
	defer c.Unlock()
	c.topics[topic] = nil

	// TODO: Send cmd...
}

// RemoveTopic is remove topics from listening.
func (c *Connection) RemoveTopic(topic string) {
	if _, ok := c.topics[topic]; !ok {
		return
	}

	c.Lock()
	defer c.Unlock()
	delete(c.topics, topic)

	// TODO: Send cmd...
}

// Topics returns all current listen topics.
func (c *Connection) Topics() []string {
	// TODO: ...
	return []string{}
}

// Close is close connection and shutdown all goroutines.
// Usually it's need to call before destroying.
func (c *Connection) Close() error {
	c.active = false
	close(c.done)

	// It can be not initialized
	if c.Connection != nil {
		return c.Connection.Close()
	}

	return nil
}

// -----------------------------------------------------------------------------

func (c *Connection) OnConnect(fn func(*Connection)) {
	c.eventOnConnect = fn
}

func (c *Connection) OnDisconnect(fn func(*Connection)) {
	c.eventOnDisconnect = fn
}

func (c *Connection) OnError(fn func(*Connection, error)) {
	c.eventOnError = fn
}

func (c *Connection) OnInfo(fn func(*Connection, string)) {
	c.eventOnInfo = fn
}

func (c *Connection) OnMessage(fn func(*Connection, []byte)) {
	c.eventOnMessage = fn
}

func (c *Connection) OnPing(fn func(*Connection, time.Time)) {
	c.eventOnPing = fn
}

func (c *Connection) OnPong(fn func(*Connection, time.Time, time.Time)) {
	c.eventOnPong = fn
}
