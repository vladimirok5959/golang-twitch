package pubsub

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

func go_pinger(c *Connection) {
	// Pinger (sender)
	go func(c *Connection) {
		for {
			select {
			case <-time.After(1 * time.Second):
				if c.active && !c.ping_sended {
					if time.Since(c.ping_start) > TwitchApiPingEach {
						if err := c.Connection.WriteMessage(
							websocket.TextMessage,
							Response{Type: "PING"}.JSON(),
						); err != nil {
							c.onError(err)
							c.active = false
							c.onDisconnect()
						} else {
							c.ping_start = time.Now()
							c.ping_sended = true
							c.onPing(c.ping_start)
						}
					}
				}
			case <-c.done:
				return
			}
		}
	}(c)

	// Pinger (handler)
	go func(c *Connection) {
		for {
			select {
			case <-time.After(1 * time.Second):
				if c.active && c.ping_sended {
					if time.Since(c.ping_start) > TwitchApiPingTimeout {
						c.onInfo(fmt.Sprintf("warning, no PONG response more than %d seconds", TwitchApiPingTimeout))
						c.active = false
						c.onDisconnect()
						c.ping_start = time.Now()
						c.ping_sended = false
						if err := c.Connection.Close(); err != nil {
							c.onError(err)
						}
					}
				}
			case <-c.done:
				return
			}
		}
	}(c)
}
