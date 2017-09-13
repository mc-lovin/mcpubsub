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

		fmt.Println(subscriber)

		publisher.Publish("got", "As1")

		//		time.Sleep(1000 * time.Millisecond)

		// we have some latency issues over here

		subscriber.UnSubscribe("got")

		/*

			subscriber1 := pubSubApi.NewSubscriber()
			subscriber1.Subscribe("got", func() {
				fmt.Println("in the callback1")
			})

			publisher.Publish("got", "As2")
		*/
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
