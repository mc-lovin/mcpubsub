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
	id int
	rw *bufio.ReadWriter
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
		id: subscriberId,
	}
	subscriberId++

	err = sendMessage(subscriberObj.rw, serverMessage{
		Id:    subscriberObj.id,
		Class: SUBSCRIBER_ADDED_MESSAGE,
	})

	if err != nil {
		log.Println("Error while creating subscriber")
	}
	return subscriberObj
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
	return err
}
