package mcpubsub

import (
	"bufio"
	"errors"
	"fmt"
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
	subscriberId                                    = 1
	subscriptionMap map[string]map[*subscriber]bool = make(map[string]map[*subscriber]bool)
)

func newSubscriber(rw *bufio.ReadWriter) subscriberApi {
	subscriberObj := &subscriber{
		rw: rw,
	}

	err := sendMessage(subscriberObj.rw, serverMessage{
		Class: SUBSCRIBER_ADDED_MESSAGE,
	})

	if err != nil {
		log.Println("Error while creating subscriber.")
	}

	message := <-channelMap[SUBSCRIBER_ADDED_MESSAGE]

	if err != nil {
		log.Println("Error while creating publisher.")
		return nil
	}

	subscriberObj.id = message.Id
	subscriberObj.callBackMap = make(map[string]Func)
	return subscriberObj
}

func (subscriberObj *subscriber) handleMessage(message serverMessage) {
	callback := subscriberObj.callBackMap[message.Topic]
	callback(message.Message)
}

func (subscriberObj *subscriber) Subscribe(topic string, callback Func) error {
	subscriberObj.callBackMap[topic] = callback
	subscriberObj.addStream(serverMessage{
		Id:    subscriberObj.id,
		Class: SUBSCRIBER_SUBSCRIBED_MESSAGE,
		Topic: topic,
	})
	return nil
}

func (subscriberObj *subscriber) UnSubscribe(topic string) error {
	subscriberObj.removeStream(serverMessage{
		Id:    subscriberObj.id,
		Class: SUBSCRIBER_UNSUBSCRIBED_MESSAGE,
		Topic: topic,
	})
	delete(subscriberObj.callBackMap, topic)
	return nil
}

func (subscriberObj *subscriber) addStream(message serverMessage) (err error) {
	topic := message.Topic
	_, ok := subscriptionMap[topic]
	if !ok {
		subscriptionMap[topic] = make(map[*subscriber]bool)
	}

	if hasElement(subscriptionMap[topic], *subscriberObj) {
		return errors.New("Subscription already exists.")
	}

	subscriptionMap[topic][subscriberObj] = true
	fmt.Println("topic ", topic, subscriberObj, subscriptionMap)
	return nil
}

func (subscriberObj *subscriber) removeStream(message serverMessage) (err error) {
	topic := message.Topic

	if !hasElement(subscriptionMap[topic], *subscriberObj) {
		return errors.New("No subscription exists.")
	}

	delete(subscriptionMap[topic], subscriberObj)
	return nil
}

func hasElement(subscriptions map[*subscriber]bool, subscriberObj subscriber) bool {
	_, ok := subscriptions[&subscriberObj]
	return ok
}
