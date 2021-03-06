package mcpubsub

import (
	"bufio"
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
	publisherId = 1
)

func newPublisher(rw *bufio.ReadWriter) publisherApi {
	publisherObj := publisher{
		rw: rw,
	}

	err := sendMessage(publisherObj.rw, serverMessage{
		Class: publisherAddedMessage,
	})
	if err != nil {
		log.Println("Error while creating publisher.")
		return nil
	}

	message := <-channelMap[publisherAddedMessage]

	if err != nil {
		log.Println("Error while creating publisher.")
		return nil
	}

	publisherObj.id = message.Id
	return publisherObj
}

func (publisherObj publisher) Publish(topic string, message string) error {
	err := sendMessage(publisherObj.rw, serverMessage{
		Id:      publisherObj.id,
		Class:   publisherPublishedMessage,
		Topic:   topic,
		Message: message,
	})
	return err
}
