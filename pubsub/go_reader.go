package pubsub

import (
	"encoding/json"
	"fmt"
	"time"
)

func go_reader(c *Connection) {
	go func(c *Connection) {
		for {
			select {
			case <-c.done:
				return
			default:
				if c.active {
					_, msg, err := c.Connection.ReadMessage()
					if err != nil {
						c.onError(err)
						c.active = false
						c.onDisconnect()

						// Wait 1 second or return immediately
						select {
						case <-time.After(time.Second):
						case <-c.done:
							return
						}
					} else {
						var answer Answer
						if err := json.Unmarshal(msg, &answer); err != nil {
							c.onError(err)
						} else {
							if answer.Type == Pong {
								ct := time.Now()
								c.onPong(c.ping_start, ct)
								c.ping_start = ct
								c.ping_sended = false
							} else if answer.Type == Reconnect {
								c.onInfo(fmt.Sprintf("warning, got %s response", Reconnect))
								c.active = false
								c.onDisconnect()
								c.ping_start = time.Now()
								c.ping_sended = false
								if err := c.Connection.Close(); err != nil {
									c.onError(err)
								}
							} else if answer.Type == Response {
								if answer.HasError() {
									c.onError(fmt.Errorf(answer.Error))
								} else {
									c.onInfo(fmt.Sprintf("type: %s, data: %#v", answer.Type, answer.Data))
								}
							} else {
								c.onMessage(msg)
							}
						}
					}
				} else {
					// Wait 1 second or return immediately
					select {
					case <-time.After(time.Second):
					case <-c.done:
						return
					}
				}
			}
		}
	}(c)
}
