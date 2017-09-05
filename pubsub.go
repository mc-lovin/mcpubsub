package main

import "fmt"

type PubSubServer struct {
}

type PubSubClient struct {
}

type callBack func()

func (pubSubServer PubSubServer) Publish(eventName string) {
	fmt.Println(eventName, pubSubServer)
}

func (pubSubClient PubSubClient) Subscribe(eventName string, callback callBack) {
	fmt.Println(eventName, pubSubClient)
	callback()
}

func (pubSubClient PubSubClient) UnSubscribe(eventName string) {
	fmt.Println(eventName, pubSubClient)
}

func throwSpoilers() {
	fmt.Println("just testing")
}

func test() {
	publisher := PubSubServer{}
	publisher.Publish("gameofthrones")

	client := PubSubClient{}
	client.Subscribe("gameofthrones", throwSpoilers)

	client.UnSubscribe("gameofthrones")
}

func main() {
	test()
}
