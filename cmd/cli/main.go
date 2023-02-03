//usr/bin/env go run "$0" "$@"; exit

package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/vladimirok5959/golang-twitch/pubsub"
)

func main() {
	ps := pubsub.New()
	defer ps.Close()

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

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	ctx := context.Background()

	go func(interrupt chan os.Signal) {
		reader := bufio.NewReader(os.Stdin)
		for {
			select {
			case <-interrupt:
				return
			default:
				fmt.Print("Enter command: ")
				cmd, _ := reader.ReadString('\n')
				cmd = strings.TrimSuffix(strings.TrimSuffix(cmd, "\r\n"), "\n")
				if strings.HasPrefix(cmd, "listen") {
					if len([]rune(cmd)) > 7 {
						param := string([]rune(cmd)[7:len([]rune(cmd))])
						fmt.Printf("Listen: (%s)\n", param)
						ps.Listen(ctx, param)
					} else {
						fmt.Printf("Parameter is not set\n")
					}
				} else if strings.HasPrefix(cmd, "unlisten") {
					if len([]rune(cmd)) > 9 {
						param := string([]rune(cmd)[9:len([]rune(cmd))])
						fmt.Printf("Unlisten: (%s)\n", param)
						ps.Unlisten(ctx, param)
					} else {
						fmt.Printf("Parameter is not set\n")
					}
				} else if strings.HasPrefix(cmd, "has") {
					if len([]rune(cmd)) > 4 {
						param := string([]rune(cmd)[4:len([]rune(cmd))])
						fmt.Printf("HasTopic: (%#v)\n", ps.HasTopic(param))
					} else {
						fmt.Printf("Parameter is not set\n")
					}
				} else if cmd == "status" {
					fmt.Printf("Status:\n")
					fmt.Printf(" - Connections count: (%d)\n", len(ps.Connections))
					fmt.Printf(" - Topics: (%#v)\n", ps.Topics())
					fmt.Printf(" - TopicsCount: (%d)\n", ps.TopicsCount())
				} else if cmd == "help" {
					fmt.Printf("Help:\n")
					fmt.Printf(" - listen <topic>\n")
					fmt.Printf(" - unlisten <topic>\n")
					fmt.Printf(" - has <topic>\n")
					fmt.Printf(" - status\n")
					fmt.Printf(" - help\n")
					fmt.Printf(" - close\n")
				} else if cmd == "close" {
					fmt.Printf("Closing...\n")
					ps.Close()
				} else {
					fmt.Printf("Unknown command\n")
				}
			}
		}
	}(interrupt)

	<-interrupt
	log.Println("Done")
}
