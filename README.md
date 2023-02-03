# golang-twitch

Twitch API PubSub interface and automatically take care of API limits. Also it will handle automatically reconnections, ping/pong and maintenance requests. Designed by PubSub guide: [https://dev.twitch.tv/docs/pubsub/](https://dev.twitch.tv/docs/pubsub/)

```go
ps := pubsub.New()

ps.OnConnect(func(c *pubsub.Connection) {
    log.Printf("OnConnect (ID: %d)\n", c.ID)
})

ps.OnDisconnect(func(c *pubsub.Connection) {
    log.Printf("OnDisconnect (ID: %d)\n", c.ID)
})

ps.OnError(func(c *pubsub.Connection, err error) {
    log.Printf("OnError (ID: %d), err: %s\n", c.ID, err)
})

ps.OnInfo(func(c *pubsub.Connection, str string) {
    log.Printf("OnInfo (ID: %d), str: %s\n", c.ID, str)
})

ps.OnMessage(func(c *pubsub.Connection, msg *pubsub.Answer) {
    log.Printf("OnMessage (ID: %d), msg: %#v\n", c.ID, msg)
})

ps.OnPing(func(c *pubsub.Connection, start time.Time) {
    log.Printf("OnPing (ID: %d), start: %d\n", c.ID, start.Unix())
})

ps.OnPong(func(c *pubsub.Connection, start, end time.Time) {
    log.Printf("OnPong (ID: %d), start: %d, end: %d\n", c.ID, start.Unix(), end.Unix())
})

ps.Listen("community-points-channel-v1", "<UserID>")

interrupt := make(chan os.Signal, 1)
signal.Notify(interrupt, os.Interrupt)
<-interrupt

ps.Close()
```

Full example here: [https://github.com/vladimirok5959/golang-twitch/blob/main/cmd/cli/main.go](https://github.com/vladimirok5959/golang-twitch/blob/main/cmd/cli/main.go)
