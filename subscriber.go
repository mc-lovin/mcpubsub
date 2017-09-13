package main

import (
	"bufio"
	"log"
)

type subscriberApi interface {
	Subscribe(topic string, callback Func) error
	UnSubscribe(topic string) error
}

type subscriber struct {
	id          int
	rw          *bufio.ReadWriter
	callBackMap map[string]Func
}

var (
	subscriberId                    = 1
	SUBSCRIBER_ADDED_MESSAGE        = "SUBSCRIBER ADDED"
	SUBSCRIBER_SUBSCRIBED_MESSAGE   = "SUBSCRIBER SUBSCRIBED"
	SUBSCRIBER_UNSUBSCRIBED_MESSAGE = "SUBSCRIBER UNSUBSCRIBED"
)

func newSubscriber(server pubSubServerApi) subscriberApi {
	rw, err := server.newRW()
	if err != nil {
		log.Println("Error while creating subscriber")
		return nil
	}
	subscriberObj := subscriber{
		rw: rw,
	}

	err = sendMessage(subscriberObj.rw, serverMessage{
		Class: SUBSCRIBER_ADDED_MESSAGE,
	})

	if err != nil {
		log.Println("Error while creating subscriber")
	}

	message, err := receiveMessage(rw)

	if err != nil {
		log.Println("Error while creating publisher")
		return nil
	}

	subscriberObj.id = message.Id
	subscriberObj.callBackMap = make(map[string]Func)

	go subscriberObj.handleMessage(rw)
	return subscriberObj
}

func (subscriberObj subscriber) handleMessage(rw *bufio.ReadWriter) {
	for {
		message, err := receiveMessage(rw)
		if err != nil {
			log.Println("error occured while subscriber recieving")
			return
		}
		callback := subscriberObj.callBackMap[message.Topic]
		callback()
		log.Print("message received in client ", message, subscriberObj.id)
	}
}

func (subscriberObj subscriber) Subscribe(topic string, callback Func) error {
	err := sendMessage(subscriberObj.rw, serverMessage{
		Id:    subscriberObj.id,
		Class: SUBSCRIBER_SUBSCRIBED_MESSAGE,
		Topic: topic,
	})

	if err != nil {
		log.Println("Error while subscribing")
	}

	subscriberObj.callBackMap[topic] = callback
	return err
}

func (subscriberObj subscriber) UnSubscribe(topic string) error {
	err := sendMessage(subscriberObj.rw, serverMessage{
		Id:    subscriberObj.id,
		Class: SUBSCRIBER_UNSUBSCRIBED_MESSAGE,
		Topic: topic,
	})

	if err != nil {
		log.Println("Error while subscribing")
	}
	delete(subscriberObj.callBackMap, topic)
	return err
}
