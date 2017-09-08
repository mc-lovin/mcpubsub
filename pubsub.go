package main

import (
	"fmt"
	"time"
)

type channel chan struct{}

type PubSubServer struct {
}

type PubSubClient struct {
	id    int
	store map[string]channel
}

// Shared Data
var (
	pubSubClientId        int                  = 1
	pubSubIdToInstanceMap map[int]PubSubClient = make(map[int]PubSubClient)
	done                                       = make(channel)
)

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
	close(channel)
	delete(pubSubClient.store, event)
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
			select {
			case channel <- struct{}{}:
			case <-done:
				return
			}
		}
	}
}

func _callOnSignalFire(buffer channel, callback callBack) {
	for {
		select {
		case _, ok := <-buffer:
			if !ok {
				return
			}
			callback()
		case <-done:
			return
		}
	}
}

func (pubSubClient PubSubClient) Subscribe(event string, callback callBack) {
	channel := pubSubClient.getOrCreateChannel(event)
	go _callOnSignalFire(channel, callback)
}

func (pubSubClient PubSubClient) UnSubscribe(event string) {
	pubSubClient.removeChannel(event)
}

func throwSpoilers() {
	fmt.Println("1got")
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

	publisher1 := createPubSubServer()

	client := createPubSubClient()

	client.Subscribe("gameofthrones", throwSpoilers)

	publisher.Publish("gameofthrones")

	client2 := createPubSubClient()
	client2.Subscribe("gameofthrones", func() {
		fmt.Println("2got")
	})

	client2.Subscribe("rickandmorty", func() {
		fmt.Println("2r&m")
	})

	publisher.Publish("gameofthrones")

	publisher.Publish("gameofthrones")

	publisher1.Publish("rickandmorty")

	client.UnSubscribe("gameofthrones")

	publisher.Publish("gameofthrones")

	hang()
}

func main() {
	defer close(done)
	test()
}

func hang() {
	// hangs the code
	time.Sleep(1000 * time.Millisecond)
}
