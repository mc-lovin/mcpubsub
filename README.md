Simple PUB-SUB implementation in go

Design
======

Subscription is topic based
Usage:


Publisher

```
publisher := PubSubServer{}
publisher.publish(topicName, context)
```

Subscriber

```
subscriber := PubSubClient{}
subscriber.subscribe(topicName, callback) 
subsrciber.unsubscribe(topicName)

```


Requriements
------------
1. publisher should be able to pass context to subscriber
2. publisher should be able to publish more than one events at ease.
3. subscriber should be able to subscribe to more than one events at ease
4. For the initial phase lets start with one publisher and multiple subscribers