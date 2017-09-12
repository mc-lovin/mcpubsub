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
	return nil
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

	var message serverMessage

	for {
		log.Println("waiting to receive")
		dec := gob.NewDecoder(rw)
		err := dec.Decode(&message)
		if err != nil {
			log.Println("Error decoding data")
			return
		}
		log.Println("Decoded message of type", message)
		handleMessage(message)
	}
}

func handleMessage(message serverMessage) {
	switch message.Class {
	case PUBLISHER_ADDED_MESSAGE:
	case PUBLISHER_PUBLISHED_MESSAGE:
	case SUBSCRIBER_ADDED_MESSAGE:
	case SUBSCRIBER_SUBSCRIBED_MESSAGE:
	case SUBSCRIBER_UNSUBSCRIBED_MESSAGE:
		return
	default:
		log.Println("Unrecognised Message")
	}
}
