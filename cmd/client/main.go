package main

import (
	"fmt"
	"log"

	"github.com/mdmdj/learn-pub-sub-starter/internal/gamelogic"
	"github.com/mdmdj/learn-pub-sub-starter/internal/pubsub"
	"github.com/mdmdj/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")
	const rabbitConnString = "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	fmt.Println("Peril game client connected to RabbitMQ!")

	fmt.Println("Getting channel...")
	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("could not get channel: %v", err)
	}

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("could not get username: %v", err)
	}
	gs := gamelogic.NewGameState(username)

	// subscribe to pause state
	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+gs.GetUsername(),
		routing.PauseKey,
		pubsub.SimpleQueueTransient,
		handlerPause(gs),
	)
	if err != nil {
		log.Fatalf("could not subscribe to pause: %v", err)
	}

	// subscribe to army moves
	err = pubsub.SubscribeJSON(
		conn,                       // connection
		routing.ExchangePerilTopic, // exchange
		routing.ArmyMovesPrefix+"."+gs.GetUsername(), // queue name
		routing.ArmyMovesPrefix+".*",                 // routing key
		pubsub.SimpleQueueTransient,                  // queue type
		handlerArmyMoves(gs),                         // handler function
	)
	if err != nil {
		log.Fatalf("could not subscribe to army moves: %v", err)
	}

	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}
		switch words[0] {
		case "move":
			am, err := gs.CommandMove(words)
			if err != nil {
				fmt.Println(err)
				continue
			}
			errPub := pubsub.PublishJSON(
				channel,                    //channel
				routing.ExchangePerilTopic, //exchange
				routing.ArmyMovesPrefix+"."+gs.GetUsername(), //key
				am, //val
			)
			if errPub != nil {
				fmt.Println("Error publishing move:", errPub)
				continue
			} else {
				fmt.Println("Published move:", am)
			}
		case "spawn":
			err = gs.CommandSpawn(words)
			if err != nil {
				fmt.Println(err)
				continue
			}
		case "status":
			gs.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			// TODO: publish n malicious logs
			fmt.Println("Spamming not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			return
		default:
			fmt.Println("unknown command")
		}
	}
}
