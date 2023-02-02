package pubsub

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

func go_reconnector(c *Connection) {
	go func(c *Connection) {
		for {
			select {
			case <-c.done:
				return
			default:
				if !c.active && len(c.topics) > 0 {
					c.onInfo(fmt.Sprintf("reconnecting to: %s", c.url.String()))
					conn, _, err := websocket.DefaultDialer.Dial(c.url.String(), nil)
					if err != nil {
						c.onError(err)

						// Wait 1 second or return immediately
						select {
						case <-time.After(time.Second):
						case <-c.done:
							return
						}
					} else {
						c.onInfo("reconnected successfully")
						c.ping_start = time.Now()
						c.ping_sended = false
						c.Connection = conn
						c.active = true
						c.onConnect()

						// Listen all topics
						c.listenTopis()
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
