/**
We will run a server at port PORT
all communication will go through that server

e.g.

Start ()

x = newPublisher()
x.publish() <- will publish to the server

y = newSubscriber()
y.subscribe() <- will subscribe to the server
y, topic should be unique pair

x and y are clients are will be given response
*/

package main

import (
	"bufio"
	"errors"
	"log"
)

type Func func()

type pubSubApi interface {
	NewPublisher() publisherApi
	NewSubscriber() subscriberApi
}

type pubSubFactory struct {
	server pubSubServerApi
	rw     *bufio.ReadWriter
}

func PubSubApi() (pubSubApi, error) {
	initChannelMap()
	server := pubSubServer{}
	rw, _ := server.newRW()
	go listener(rw)
	go broadCaster()
	sendMessage(rw, serverMessage{
		Class: ADD_CLIENT_MESSAGE,
	})
	pubSubObj := pubSubFactory{
		server: server,
		rw:     rw,
	}
	return pubSubObj, nil
}

func listener(rw *bufio.ReadWriter) {
	for {
		message, err := receiveMessage(rw)
		if err != nil {
			log.Println("not able to receive messages atm")
			return
		}
		log.Println("in listener", message)
		channelMap[message.Class] <- message
	}
}

func broadCaster() {
	for {
		message := <-channelMap[PUBLISHER_PUBLISHED_MESSAGE]
		broadCastMessage(message)
	}
}

func broadCastMessage(message serverMessage) (err error) {

	if message.Id == 0 {
		return errors.New("Invalid Message, Id is null")
	}

	streams, ok := subscriptionMap[message.Topic]
	if !ok {
		log.Println("No subscribers yet for ", message)
		return nil
	}

	for subscriber := range streams {
		log.Println("Sending ", subscriber, message.Message)
		_, ok = subscriber.callBackMap[message.Topic]
		if !ok {
			continue
		}
		subscriber.handleMessage(message)
	}
	return nil
}

func PubSubApiServerStart() (pubSubApi, error) {
	pubSubObj := pubSubFactory{
		server: pubSubServer{},
	}
	err := pubSubObj.server.start()
	if err != nil {
		return nil, err
	}
	return pubSubObj, nil
}

func (pubSubObj pubSubFactory) NewPublisher() publisherApi {
	return newPublisher(pubSubObj.rw)
}

func (pubSubObj pubSubFactory) NewSubscriber() subscriberApi {
	return newSubscriber(pubSubObj.rw)
}
