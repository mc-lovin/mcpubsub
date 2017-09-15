package main

import (
	"flag"
	"fmt"
	"log"
	"time"
	"github.com/mc-lovin/mcpubsub"
)

func main() {
	runServer := flag.Bool("runserver", false, "run a server")
	flag.Parse()

	if *runServer {
		_, err := mcpubsub.PubSubApiServerStart()
		if err != nil {
			log.Println("Error in starting server.")
		}
	} else {
		pubSubApi, err := mcpubsub.PubSubApi()
		if err != nil {
			log.Println("Error in listening.")
			return
		}
		publisher := pubSubApi.NewPublisher()
		subscriber := pubSubApi.NewSubscriber()

		subscriber.Subscribe("got", func(message string) {
			fmt.Println("Game of thrones: ", message)
		})
		publisher.Publish("got", "All hail Lord Baelish")

		subscriber.Subscribe("r&m", func(message string) {
			fmt.Println("Rick and Morty", message)
		})

		time.Sleep(time.Millisecond)
		subscriber.UnSubscribe("got")
		publisher.Publish("r&m", "Dabba wabba doo")
		publisher.Publish("got", "All hail Lord Baelish")
	}
	hang()
}

func hang() {
	time.Sleep(10000 * time.Millisecond)
}
