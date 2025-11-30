package main

import (
	"fmt"

	"github.com/mdmdj/learn-pub-sub-starter/internal/gamelogic"
	"github.com/mdmdj/learn-pub-sub-starter/internal/pubsub"
	"github.com/mdmdj/learn-pub-sub-starter/internal/routing"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	connectionString := "amqp://guest:guest@localhost:5672/"
	amqpConnection, err := amqp.Dial(connectionString)
	if err != nil {
		fmt.Println("Error connecting to RabbitMQ: ", err)
		panic("Error connecting to RabbitMQ")
	}
	defer amqpConnection.Close()

	fmt.Println("Connected", amqpConnection.Properties)

	mqChannel, err := amqpConnection.Channel()
	if err != nil {
		fmt.Println("Error creating channel: ", err)
		panic("Error creating channel")
	}
	defer mqChannel.Close()

	logChannel, logQueue, err := bindToLogs(amqpConnection)
	if err != nil {
		fmt.Println("Error in bindToLogs: ", err)
		panic("Error in bindToLogs")
	}
	fmt.Println(logChannel, logQueue)

	gamelogic.PrintServerHelp()

	replRunning := true

	for replRunning {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}

		fmt.Println("Input: ", words, len(words))

		switch words[0] {
		case "pause":
			fmt.Println("Sending pause message from server...")
			err := sendPause(mqChannel)
			if err != nil {
				fmt.Printf("Error sending pause: %v", err)
			}
		case "resume":
			fmt.Println("Sending resume message from server...")
			err := sendResume(mqChannel)
			if err != nil {
				fmt.Printf("Error sending resume: %v", err)
			}
		case "quit":
			fmt.Println("Quitting server...")
			replRunning = false
		case "help":
			gamelogic.PrintServerHelp()
		default:
			fmt.Println("Unknown command: ", words[0])
			gamelogic.PrintServerHelp()
		}
	}
}

func sendPause(channel *amqp.Channel) error {
	return pubsub.PublishJSON(channel,
		routing.ExchangePerilDirect,
		routing.PauseKey,
		routing.PlayingState{IsPaused: true})
}

func sendResume(channel *amqp.Channel) error {
	return pubsub.PublishJSON(channel,
		routing.ExchangePerilDirect,
		routing.PauseKey,
		routing.PlayingState{IsPaused: false})
}

func bindToLogs(amqpConnection *amqp.Connection) (
	channel *amqp.Channel,
	queue amqp.Queue,
	err error) {
	channel, queue, err = pubsub.DeclareAndBind(
		amqpConnection,
		routing.ExchangePerilTopic, routing.GameLogSlug,
		routing.GameLogSlug+".*", pubsub.SimpleQueueDurable)
	return
}
