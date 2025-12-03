package main

import (
	"fmt"

	"github.com/mdmdj/learn-pub-sub-starter/internal/gamelogic"
	"github.com/mdmdj/learn-pub-sub-starter/internal/pubsub"
	"github.com/mdmdj/learn-pub-sub-starter/internal/routing"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) pubsub.AckType {
	return func(ps routing.PlayingState) (at pubsub.AckType) {
		defer fmt.Print("> ")
		fmt.Println("Handling Pause: ", ps)
		gs.HandlePause(ps)
		fmt.Println("Ack")
		return pubsub.Ack
	}
}

func handlerArmyMoves(gs *gamelogic.GameState) func(gamelogic.ArmyMove) pubsub.AckType {
	return func(am gamelogic.ArmyMove) (at pubsub.AckType) {
		defer fmt.Print("> ")
		fmt.Println("Handling Move: ", am)
		mc := gs.HandleMove(am)
		switch mc {
		case gamelogic.MoveOutComeSafe:
			at = pubsub.Ack
			fmt.Println("Ack")
			return
		case gamelogic.MoveOutcomeMakeWar:
			at = pubsub.Ack
			fmt.Println("Ack")
			return
		case gamelogic.MoveOutcomeSamePlayer:
			at = pubsub.NackDiscard
			fmt.Println("NackDiscard")
			return
		default:
			at = pubsub.NackDiscard
			fmt.Println("NackDiscard")
			return
		}
	}
}
