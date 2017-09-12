/**
We will run a server at port PORT
all communication will go through that server

e.g.

Start ()

x = newPublisher()
x.publish() <- will publish to the server

y = newSubscriber()
y.subscribe() <- will subscribe to the server
y, topic should be unique pair

x and y are clients are will be given response
*/

package main

type Func func()

type pubSubApi interface {
	NewPublisher() publisherApi
	NewSubscriber() subscriberApi
}

type pubSubFactory struct {
	server pubSubServerApi
}

func PubSubApi() (pubSubApi, error) {
	pubSubObj := pubSubFactory{
		server: pubSubServer{},
	}
	return pubSubObj, nil
}

func PubSubApiServerStart() (pubSubApi, error) {
	pubSubObj := pubSubFactory{
		server: pubSubServer{},
	}
	err := pubSubObj.server.start()
	if err != nil {
		return nil, err
	}
	return pubSubObj, nil
}

func (pubSubObj pubSubFactory) NewPublisher() publisherApi {
	return newPublisher(pubSubObj.server)
}

func (pubSubObj pubSubFactory) NewSubscriber() subscriberApi {
	return newSubscriber(pubSubObj.server)
}
