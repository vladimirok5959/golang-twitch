// Package implements Twitch API PubSub and automatically take care of API
// limit. Also it will handle automatically reconnections, ping/pong and
// maintenance requests.
package pubsub

import (
	"fmt"
	"net/url"
	"strings"
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
	URL         url.URL
	Connections map[int64]*Connection
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
	p := PubSub{URL: url}
	return &p
}

// Listen is adding topics for listening. It take care of API limits.
// New TCP connection will be created for every 50 topics.
//
// https://dev.twitch.tv/docs/pubsub/#connection-management
func (p *PubSub) Listen(topic string, params ...interface{}) {
	// TODO: ...
}

// Unlisten is remove topics from listening. It take care of API limits too.
// Connection count will automatically decrease of needs.
//
// https://dev.twitch.tv/docs/pubsub/#connection-management
func (p *PubSub) Unlisten(topic string, params ...interface{}) {
	// TODO: ...
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
