package mcpubsub

const (
	addClientMessage              = "CLIENT ADDED"
	publisherAddedMessage         = "PUBLISHER ADDED"
	publisherPublishedMessage     = "PUBLISHER PUBLISHED"
	subscriberAddedMessage        = "SUBSCRIBER ADDED"
	subscriberSubscribedMessage   = "SUBSCRIBER SUBSCRIBED"
	subscriberUnsubscribedMessage = "SUBSCRIBER UNSUBSCRIBED"
	chanLen                       = 1
)

var (
	channelMap map[string]chan serverMessage = make(map[string]chan serverMessage)
)

func initChannelMap() {
	channelMap[publisherAddedMessage] = make(chan serverMessage, chanLen)
	channelMap[publisherPublishedMessage] = make(chan serverMessage, chanLen)
	channelMap[subscriberAddedMessage] = make(chan serverMessage, chanLen)
	channelMap[subscriberSubscribedMessage] = make(chan serverMessage, chanLen)
	channelMap[subscriberUnsubscribedMessage] = make(chan serverMessage, chanLen)
}
