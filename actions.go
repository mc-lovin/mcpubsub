package main

const (
	ADD_CLIENT_MESSAGE              = "CLIENT ADDED"
	PUBLISHER_ADDED_MESSAGE         = "PUBLISHER ADDED"
	PUBLISHER_PUBLISHED_MESSAGE     = "PUBLISHER PUBLISHED"
	SUBSCRIBER_ADDED_MESSAGE        = "SUBSCRIBER ADDED"
	SUBSCRIBER_SUBSCRIBED_MESSAGE   = "SUBSCRIBER SUBSCRIBED"
	SUBSCRIBER_UNSUBSCRIBED_MESSAGE = "SUBSCRIBER UNSUBSCRIBED"
	chanLen                         = 1
)

var (
	channelMap map[string]chan serverMessage = make(map[string]chan serverMessage)
)

func initChannelMap() {
	channelMap[PUBLISHER_ADDED_MESSAGE] = make(chan serverMessage, chanLen)
	channelMap[PUBLISHER_PUBLISHED_MESSAGE] = make(chan serverMessage, chanLen)
	channelMap[SUBSCRIBER_ADDED_MESSAGE] = make(chan serverMessage, chanLen)
	channelMap[SUBSCRIBER_SUBSCRIBED_MESSAGE] = make(chan serverMessage, chanLen)
	channelMap[SUBSCRIBER_UNSUBSCRIBED_MESSAGE] = make(chan serverMessage, chanLen)
}
