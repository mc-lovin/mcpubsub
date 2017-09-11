package main

import (
	"bufio"
)

type publisherApi interface {
	publish(topic string, message string) error
}

var (
	publisherId             = 1
	PUBLISHER_ADDED_MESSAGE = "PUBLISHER ADDED"
)

func newPublisher(server pubSubServerApi) publisherApi {
	publisherObj := publisher{
		rw: server.getReadWriter(),
		id: publisherId,
	}
	publisherId++
	sendMessage(publisherObj.rw, serverMessage{
		id:    publisherObj.id,
		class: PUBLISHER_ADDED_MESSAGE,
	})
	return publisherObj
}

type publisher struct {
	id int
	rw *bufio.ReadWriter
}

func (publisherObj publisher) write(topic string, message string) error {
	sendMessage(publisherObj.rw, serverMessage{
		id:    publisherObj.id,
		class: PUBLISHER_ADDED_MESSAGE,
	})
	return nil
}

func (publisherObj publisher) publish(topic string, message string) error {
	return nil
}
