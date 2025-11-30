package main

import (
	"fmt"

	"github.com/mdmdj/learn-pub-sub-starter/internal/gamelogic"
	"github.com/mdmdj/learn-pub-sub-starter/internal/pubsub"
	"github.com/mdmdj/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")

	connectionString := "amqp://guest:guest@localhost:5672/"
	amqpConnection, err := amqp.Dial(connectionString)
	if err != nil {
		fmt.Println("Error connecting to RabbitMQ: ", err)
		panic("Error connecting to RabbitMQ")
	}
	defer amqpConnection.Close()

	fmt.Println("Connected", amqpConnection.Properties)

	userName, err := gamelogic.ClientWelcome()
	if err != nil {
		fmt.Println("Error in gamelogic.ClientWelcome: ", err)
		panic("Error in gamelogic.ClientWelcome")
	}

	fmt.Println(userName)
	fmt.Println(err)

	bindToPause(amqpConnection, userName)
	gameState := gamelogic.NewGameState(userName)

	gamelogic.PrintClientHelp()
	replRunning := true

	for replRunning {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}

		fmt.Println("Input: ", words, len(words))

		switch words[0] {
		case "quit":
			fmt.Println("Quitting client...")
			gamelogic.PrintQuit()
			replRunning = false
		case "help":
			gamelogic.PrintClientHelp()
		case "spawn":
			err := gameState.CommandSpawn(words)
			fmt.Println(err)
		case "move":
			move, err := gameState.CommandMove(words)
			fmt.Println(move, err)
		case "status":
			gameState.CommandStatus()
		case "spam":
			fmt.Println("Spamming not allowed yet!")
		default:
			fmt.Println("Unknown command: ", words[0])
			gamelogic.PrintClientHelp()
		}
	}

}

func bindToPause(amqpConnection *amqp.Connection, userName string) (err error) {
	pauseUser := routing.PauseKey + "." + userName
	channel, queue, err := pubsub.DeclareAndBind(amqpConnection,
		routing.ExchangePerilDirect, pauseUser,
		routing.PauseKey, pubsub.Transient)
	fmt.Println(channel)
	fmt.Println(queue)
	if err != nil {
		fmt.Println("Error in pubsub.DeclareAndBind: ", err)
		//panic("Error in pubsub.DeclareAndBind")
	}
	return
}
