package main

import (
	"bufio"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"net"
)

var (
	PORT = ":14610"
	// Tells which streams to publish
	subscriptionMap map[string]map[int]bool   = make(map[string]map[int]bool)
	subStreamMap    map[int]*bufio.ReadWriter = make(map[int]*bufio.ReadWriter)
	pubStreamMap    map[int]*bufio.ReadWriter = make(map[int]*bufio.ReadWriter)
)

type pubSubServerApi interface {
	newRW() (*bufio.ReadWriter, error)
	start() error
}

type pubSubServer struct {
}

type serverMessage struct {
	Id      int
	Class   string
	Topic   string
	Message string
}

func sendMessage(rw *bufio.ReadWriter, message serverMessage) error {
	enc := gob.NewEncoder(rw)
	err := enc.Encode(message)
	if err != nil {
		log.Println("error while sending message", message, err)
		return err
	}
	rw.Flush()
	log.Println("sent message")

	return nil
}

func receiveMessage(rw *bufio.ReadWriter) (message serverMessage, err error) {
	log.Println("waiting to receive")
	dec := gob.NewDecoder(rw)
	err = dec.Decode(&message)
	if err != nil {
		log.Println("Error decoding data")
		return message, err
	}
	log.Println("Decoded message of type", message)
	return message, nil
}

func (server pubSubServer) newRW() (*bufio.ReadWriter, error) {
	conn, err := net.Dial("tcp", "localhost"+PORT)
	if err != nil {
		return nil, err
	}
	fmt.Println("New client added")
	return bufio.NewReadWriter(bufio.NewReader(conn),
		bufio.NewWriter(conn)), nil
}

func (server pubSubServer) start() error {
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		return errors.New("Not able to accept connections atm")
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("error accepting connection")
			continue
		}
		log.Println("connection accepted")

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	for {
		message, err := receiveMessage(rw)
		if err != nil {
			fmt.Println("returning from handleconnection")
			return
		}
		handleRequestMessage(message, rw)
	}
}

func handleRequestMessage(message serverMessage, rw *bufio.ReadWriter) {
	var err error = nil
	switch message.Class {
	case PUBLISHER_ADDED_MESSAGE:
		message.Id = publisherId
		pubStreamMap[publisherId] = rw
		publisherId++
		err = sendMessage(rw, message)
	case PUBLISHER_PUBLISHED_MESSAGE:
		broadCastMessage(message)
	case SUBSCRIBER_ADDED_MESSAGE:
		message.Id = subscriberId
		subStreamMap[subscriberId] = rw
		subscriberId++
		err = sendMessage(rw, message)
	case SUBSCRIBER_SUBSCRIBED_MESSAGE:
		addStream(message)
	case SUBSCRIBER_UNSUBSCRIBED_MESSAGE:
		removeStream(message)
	default:
		log.Println("Unrecognised Message")
	}

	if err != nil {
		log.Println("Problem dealing with serverMessage", message)
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

	for subId := range streams {
		log.Println("Sending ", subId, message.Message)
		sendMessage(subStreamMap[subId], message)
	}

	return nil

}

func addStream(message serverMessage) (err error) {
	topic := message.Topic
	id := message.Id
	_, ok := subscriptionMap[topic]
	if !ok {
		subscriptionMap[topic] = make(map[int]bool)
	}

	if hasElement(subscriptionMap[topic], id) {
		return errors.New("Subscription already exists")
	}

	subscriptionMap[topic][id] = true
	return nil
}

func removeStream(message serverMessage) (err error) {
	topic := message.Topic
	id := message.Id

	if !hasElement(subscriptionMap[topic], id) {
		return errors.New("No subscription exists")
	}

	delete(subscriptionMap[topic], id)
	return nil
}

func hasElement(subscriptions map[int]bool, id int) bool {
	_, ok := subscriptions[id]
	return ok
}
