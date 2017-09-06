package main

import (
	"fmt"
	"time"
)

type channel chan int

type PubSubServer struct {
}

type PubSubClient struct {
	id    int
	store map[string]channel
}

// Shared Data
var pubSubClientId int = 1
var pubSubIdToInstanceMap map[int]PubSubClient = make(map[int]PubSubClient)

// -----------

func (pubSubClient *PubSubClient) channelExists(event string) bool {
	_, ok := pubSubClient.store[event]
	return ok
}

func (pubSubClient *PubSubClient) getOrCreateChannel(event string) channel {
	if !pubSubClient.channelExists(event) {
		ret := make(channel)
		pubSubClient.store[event] = ret
	}
	return pubSubClient.store[event]
}

func (pubSubClient *PubSubClient) removeChannel(event string) {
	if !pubSubClient.channelExists(event) {
		return
	}
	channel := pubSubClient.store[event]
	delete(pubSubClient.store, event)
	close(channel)
}

func (pubSubClient *PubSubClient) isInitialised() bool {
	return pubSubClient.store != nil
}

func (pubSubClient *PubSubClient) init() {
	if pubSubClient.isInitialised() {
		return
	}
	pubSubClient.id = pubSubClientId
	pubSubClientId = pubSubClientId + 1
	pubSubClient.store = make(map[string]channel)
	pubSubIdToInstanceMap[pubSubClient.id] = *pubSubClient
}

// todo callBack should accept variables eventually
type callBack func()

func (pubSubServer PubSubServer) Publish(event string) {
	for id := 1; id < pubSubClientId; id += 1 {
		pubSubClient := pubSubIdToInstanceMap[id]
		if pubSubClient.channelExists(event) {
			channel := pubSubClient.getOrCreateChannel(event)
			go func() {
				channel <- 1
			}()
		}
	}
}

func (pubSubClient PubSubClient) Subscribe(event string, callback callBack) {

	channel := pubSubClient.getOrCreateChannel(event)
	go func() {
		for {
			_, ok := <-channel
			if !ok {
				continue
			}
			callback()
		}
	}()
}

func (pubSubClient PubSubClient) UnSubscribe(event string) {
	pubSubClient.removeChannel(event)
}

func throwSpoilers() {
	fmt.Println("just testing")
}

func createPubSubClient() PubSubClient {
	pubSubClient := PubSubClient{}
	pubSubClient.init()
	return pubSubClient
}

func createPubSubServer() PubSubServer {
	return PubSubServer{}
}

func test() {

	publisher := createPubSubServer()

	client := createPubSubClient()

	client.Subscribe("gameofthrones", throwSpoilers)

	time.Sleep(1000 * time.Millisecond)

	publisher.Publish("gameofthrones")

	time.Sleep(2000 * time.Millisecond)

	client.UnSubscribe("gameofthrones")

	time.Sleep(2000 * time.Millisecond)

	publisher.Publish("gameofthrones")

	hang()
}

func main() {
	test()
}

func hang() {
	// hangs the code
	time.Sleep(10000 * time.Millisecond)
}
