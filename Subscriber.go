package main

type subscriberApi interface {
}

func newSubscriber(pubSubServerApi) subscriberApi {
	return subscriber{}
}

type subscriber struct {
	id int
}

func (subscriberObj subscriber) Subscribe(message string) error {
	return nil
}
