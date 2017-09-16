Simple PUB-SUB implementation in go

Usage
=====

Installation
`go get -v github.com/mc-lovin/mcpubsub`

Import `github.com/mc-lovin/mcpubsub`

Run a server using
```
mcpubsub.PubSubApiServerStart()
```

Run a pub sub client using
```
pubSubApi := mcpubsub.PubSubApi()
```

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

`callBack` will accept one `string` as an argument which is going to be
the `message`

Once connected all the clients would be able to publish messages across the network

See `test/pubsub.go` for a full example