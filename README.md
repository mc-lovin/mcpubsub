Simple PUB-SUB implementation in go

Usage
=====

Installation
`go get -v github.com/mc-lovin/mcpubsub`

Import `github.com/mc-lovin/mcpubsub`


Run a server using
`mcpubsub.PubSubApiServerStart()`

Run a pub sub client using
`pubSub := mcpubsub.PubSubApi()`

Publisher

```
publisher := pubSubApi.NewPublisher()
publisher.publish(topicName)
```

Subscriber

```
subscriber := pubSubApi.NewSubscriber()
subscriber.subscribe(topicName, callback)
subsrciber.unsubscribe(topicName)

```

callBack has to be of the type `func(string) {}`

Once connected all the clients would be able to publish messages across the network

See `test/pubsub.go` for more understanding