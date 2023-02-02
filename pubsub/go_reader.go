package pubsub

import (
	"encoding/json"
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
						var resp Response
						if err := json.Unmarshal(msg, &resp); err != nil {
							c.onError(err)
						} else {
							if resp.Type == "PONG" {
								ct := time.Now()
								c.onPong(c.ping_start, ct)
								c.ping_start = ct
								c.ping_sended = false
							} else if resp.Type == "RECONNECT" {
								c.onInfo("warning, got RECONNECT response")
								c.active = false
								c.onDisconnect()
								c.ping_start = time.Now()
								c.ping_sended = false
								if err := c.Connection.Close(); err != nil {
									c.onError(err)
								}
							} else if resp.Type == "RESPONSE" {
								// TODO: {"type":"RESPONSE","error":"","nonce":""}
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
