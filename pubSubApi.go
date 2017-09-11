/*
We will run a server at port PORT
all communication will go through that server

e.g.

Start ()

x = newPublisher()
x.publish() <- will publish to the server

y = newSubscriber()
y.subscribe() <- will subscribe to the server

x and y are clients are will be given response
*/

package main

type Func func()

type pubSubApi interface {
	NewPublisher() publisherApi
	NewSubscriber(topic string) subscriberApi
	start()
}

var (
	apiInstance pubSubApi
)

func PubSubApi() pubSubApi {
	if apiInstance != nil {
		return apiInstance
	}
	apiInstance = pubSubFactory{}
	apiInstance.start()
	return apiInstance
}

type pubSubFactory struct {
	server pubSubServerApi
}

func (pubSubFactory pubSubFactory) start() {
}

func (pubSubFactory pubSubFactory) NewPublisher() publisherApi {
	return newPublisher(pubSubFactory.server)
}

func (pubSubFactory pubSubFactory) NewSubscriber(topic string) subscriberApi {
	return newSubscriber(pubSubFactory.server)
}
