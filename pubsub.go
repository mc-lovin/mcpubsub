package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	//	defer close(done)
	//	test()

	connect := flag.String("connect", "", "IP addr")

	flag.Parse()

	if *connect != "" {
		pubSubApi, _ := PubSubApi()
		publisher := pubSubApi.NewPublisher()
		// Print publisher and see the log, interesting right.

		subscriber := pubSubApi.NewSubscriber()
		subscriber.Subscribe("got", func() {
			fmt.Println("in the callback")
		})

		publisher.Publish("got", "As")

		time.Sleep(100 * time.Millisecond)

		subscriber.UnSubscribe("got")

		subscriber1 := pubSubApi.NewSubscriber()
		subscriber1.Subscribe("got", func() {
			fmt.Println("in the callback1")
		})

		publisher.Publish("got", "As")

		fmt.Println("--->", publisher, subscriber)

	} else {
		PubSubApiServerStart()
	}

	hang()

}

func hang() {
	// hangs the code
	time.Sleep(10000 * time.Millisecond)
}
