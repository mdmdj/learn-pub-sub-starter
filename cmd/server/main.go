package main

import (
	"fmt"
	"github.com/mdmdj/learn-pub-sub-starter/internal/pubsub"
	"github.com/mdmdj/learn-pub-sub-starter/internal/routing"
	"os"
	"os/signal"
	"syscall"

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

	pubsub.PublishJSON(mqChannel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})

	// Listen for CTRL+C to close the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	interruptSignal := <-c
	fmt.Println(interruptSignal)
	fmt.Printf("Got signal: %v\n", interruptSignal)
	fmt.Println("Closing server")
}
