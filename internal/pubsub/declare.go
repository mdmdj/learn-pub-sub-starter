package pubsub

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType int

const (
	Durable SimpleQueueType = iota
	Transient
)

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
) (channel *amqp.Channel, queue amqp.Queue, err error) {
	channel, err = conn.Channel()
	if err != nil {
		return
	}

	// prepare queue attributes based on the queue type
	isDurable := (bool)(queueType == Durable)
	isAutoDelete := (bool)(queueType == Transient)
	isExclusive := (bool)(queueType == Transient)

	isNoWait := false
	queueArgs := (amqp.Table)(nil)

	queue, err = channel.QueueDeclare(
		queueName,
		isDurable,
		isAutoDelete,
		isExclusive,
		isNoWait,
		queueArgs)

	if err != nil {
		return
	}

	return
}
