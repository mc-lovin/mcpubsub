package main

import "bufio"

var (
	PORT = 14610
)

type pubSubServerApi interface {
	getReadWriter() *bufio.ReadWriter
}

type pubSubServer struct {
}

type serverMessage struct {
	id      int
	class   string
	topic   string
	message string
}

func sendMessage(rw *bufio.ReadWriter, message serverMessage) error {
	return nil
}
