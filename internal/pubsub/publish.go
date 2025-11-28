package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) (err error) {
	valJson, err := json.Marshal(val)
	if err != nil {
		fmt.Println("Error in internal/pubsub/PublishJSON while marshalling val")
		fmt.Println(err)
		return
	}

	err = ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        valJson,
	})

	if err != nil {
		fmt.Println("Error in internal/pubsub/PublishJSON while attempting to publish")
		fmt.Println(err)
		return
	}

	return
}
