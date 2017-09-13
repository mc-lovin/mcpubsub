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
	subStreamMap map[int]*bufio.ReadWriter = make(map[int]*bufio.ReadWriter)
	pubStreamMap map[int]*bufio.ReadWriter = make(map[int]*bufio.ReadWriter)
	streams      []*bufio.ReadWriter       = make([]*bufio.ReadWriter, 0)
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
	rw := bufio.NewReadWriter(bufio.NewReader(conn),
		bufio.NewWriter(conn))
	return rw, nil
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
	case ADD_CLIENT_MESSAGE:
		streams = append(streams, rw)
		fmt.Println("added streams ", streams)
	case PUBLISHER_ADDED_MESSAGE:
		message.Id = publisherId
		pubStreamMap[publisherId] = rw
		publisherId++
		err = sendMessage(rw, message)
	case PUBLISHER_PUBLISHED_MESSAGE:
		log.Println("publisher published ", streams)
		for idx, stream := range streams {
			log.Println(idx, stream)
			sendMessage(stream, message)
		}
	case SUBSCRIBER_ADDED_MESSAGE:
		message.Id = subscriberId
		subStreamMap[subscriberId] = rw
		subscriberId++
		err = sendMessage(rw, message)
	case SUBSCRIBER_SUBSCRIBED_MESSAGE:
		return
	case SUBSCRIBER_UNSUBSCRIBED_MESSAGE:
		return
	default:
		log.Println("Unrecognised Message")
	}

	if err != nil {
		log.Println("Problem dealing with serverMessage", message)
	}
}
