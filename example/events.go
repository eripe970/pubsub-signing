package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	signing "github.com/eripe970/pubsub-signing"
	"log"
)

type Event interface {
	Type() string
}

type UserChangedNameEvent struct {
	Name string
}

func (n UserChangedNameEvent) Type() string {
	return "user.name.changed"
}

func PublishEvent(projectId, topic, secret string, event Event) error {
	client, err := pubsub.NewClient(context.Background(), projectId)
	if err != nil {
		log.Fatal(err)
	}

	t := client.Topic(topic)

	data, err := json.Marshal(event)

	if err != nil {
		return fmt.Errorf("failed to serialize message: %s", err.Error())
	}

	message := pubsub.Message{
		Data: data,
		Attributes: map[string]string{
			"event-type": event.Type(),
		},
	}

	err = signing.SignMessage(&message, secret)

	if err != nil {
		return err
	}

	result := t.Publish(context.Background(), &message)

	id, err := result.Get(context.Background())
	if err != nil {
		return err
	}

	log.Printf("published message %v\n", id)
	return nil
}
