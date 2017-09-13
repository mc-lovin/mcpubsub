package main

import (
	"bufio"
	"fmt"
	"log"
)

type publisherApi interface {
	Publish(topic string, message string) error
}

type publisher struct {
	id int
	rw *bufio.ReadWriter
}

var (
	publisherId                 = 1
	PUBLISHER_ADDED_MESSAGE     = "PUBLISHER ADDED"
	PUBLISHER_PUBLISHED_MESSAGE = "PUBLISHER PUBLISHED"
)

func newPublisher(server pubSubServerApi) publisherApi {
	rw, err := server.newRW()
	if err != nil {
		log.Println("Error while creating publisher")
		return nil
	}
	publisherObj := publisher{
		rw: rw,
	}

	err = sendMessage(publisherObj.rw, serverMessage{
		Class: PUBLISHER_ADDED_MESSAGE,
	})
	if err != nil {
		log.Println("Error while creating publisher")
		return nil
	}

	message, err := receiveMessage(rw)

	if err != nil {
		log.Println("Error while creating publisher")
		return nil
	}

	publisherObj.id = message.Id
	fmt.Println(publisherObj)
	return publisherObj
}

func (publisherObj publisher) Publish(topic string, message string) error {
	err := sendMessage(publisherObj.rw, serverMessage{
		Id:      publisherObj.id,
		Class:   PUBLISHER_PUBLISHED_MESSAGE,
		Topic:   topic,
		Message: message,
	})
	return err
}
