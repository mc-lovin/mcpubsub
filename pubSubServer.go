package main

import (
	"bufio"
	"encoding/gob"
	"errors"
	"log"
	"net"
)

var (
	PORT                        = ":14610"
	streams []*bufio.ReadWriter = make([]*bufio.ReadWriter, 0)
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
	err := enc.Encode(&message)
	if err != nil {
		log.Println("Error while sending message.", message, err)
		return err
	}
	rw.Flush()
	log.Println("Sent message.", message)
	return nil
}

func receiveMessage(rw *bufio.ReadWriter) (message serverMessage, err error) {
	log.Println("Waiting to receive message.")
	dec := gob.NewDecoder(rw)
	err = dec.Decode(&message)
	if err != nil {
		log.Println("Error decoding data", err)
		return message, err
	}
	log.Println("Received message", message)
	return message, nil
}

func (server pubSubServer) newRW() (*bufio.ReadWriter, error) {
	conn, err := net.Dial("tcp", "localhost"+PORT)
	if err != nil {
		return nil, err
	}
	log.Println("New client added.")
	rw := bufio.NewReadWriter(bufio.NewReader(conn),
		bufio.NewWriter(conn))
	return rw, nil
}

func (server pubSubServer) start() error {
	listener, err := net.Listen("tcp", PORT)
	if err != nil {
		return errors.New("Not able to accept connections.")
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection", err)
			continue
		}
		log.Println("Connection Accepted.")

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	for {
		message, err := receiveMessage(rw)
		if err != nil {
			log.Println("error in handleConnection", err)
			return
		}
		handleRequestMessage(message, rw)
	}
}

func handleRequestMessage(message serverMessage, rw *bufio.ReadWriter) {
	var err error = nil
	log.Println("Handling message", message)
	switch message.Class {
	case ADD_CLIENT_MESSAGE:
		streams = append(streams, rw)
	case PUBLISHER_ADDED_MESSAGE:
		message.Id = publisherId
		publisherId++
		err = sendMessage(rw, message)
	case PUBLISHER_PUBLISHED_MESSAGE:
		for idx, stream := range streams {
			log.Println(idx, stream)

			err = sendMessage(stream, message)
			if err != nil {
				log.Println("Bad stream can't write.")
			}
		}
	case SUBSCRIBER_ADDED_MESSAGE:
		message.Id = subscriberId
		subscriberId++
		err = sendMessage(rw, message)
	case SUBSCRIBER_SUBSCRIBED_MESSAGE:
		return
	case SUBSCRIBER_UNSUBSCRIBED_MESSAGE:
		return
	default:
		log.Println("Unrecognised message.")
	}

	if err != nil {
		log.Println("Problem dealing with serverMessage.", message)
	}
}
