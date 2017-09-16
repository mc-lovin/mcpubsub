package mcpubsub

import (
	"bufio"
	"errors"
	"log"
)

type Func func(string)

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
	rw, err := server.newRW()
	if err != nil {
		log.Println("Error in starting service", err)
		return nil, err
	}

	go listener(rw)
	go broadCaster()

	sendMessage(rw, serverMessage{
		Class: addClientMessage,
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
			log.Println("Not able to receive messages.", err)
			return
		}
		channelMap[message.Class] <- message
	}
}

func broadCaster() {
	for {
		message := <-channelMap[publisherPublishedMessage]
		broadCastMessage(message)
	}
}

func broadCastMessage(message serverMessage) (err error) {

	log.Println("Broadcasting message", message)

	if message.Id == 0 {
		return errors.New("Invalid Message, Id is null.")
	}

	streams, ok := subscriptionMap[message.Topic]
	log.Println(message, streams, ok)
	if !ok {
		return nil
	}

	for subscriber := range streams {
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
