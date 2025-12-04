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
		fmt.Println(cs.GetUserName(), "Handling Move: ", am)
		mc := gs.HandleMove(am)
		switch mc {
		case gamelogic.MoveOutComeSafe:
			at = pubsub.Ack
			fmt.Println("Ack")
			return
		case gamelogic.MoveOutcomeMakeWar:
			err := pubsub.PublishJSON(
				cs.Channel,                 //channel
				routing.ExchangePerilTopic, // exchange
				routing.WarRecognitionsPrefix+"."+cs.GetUserName(), // key
				gamelogic.RecognitionOfWar{
					Attacker: am.Player,
					Defender: gs.GetPlayerSnap(),
				},
			)
			if err != nil {
				fmt.Println("Error publishing war:", err)
				at = pubsub.NackRequeue
				fmt.Println("NackRequeue")
			}
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

func handlerWar(gs *gamelogic.GameState) func(gamelogic.RecognitionOfWar) pubsub.AckType {
	return func(rw gamelogic.RecognitionOfWar) (at pubsub.AckType) {
		defer fmt.Print("> ")
		fmt.Println(cs.GetUserName(), "Handling War: ", rw)
		oc, _, _ := gs.HandleWar(rw)
		switch oc {
		case gamelogic.WarOutcomeNotInvolved:
			at = pubsub.NackRequeue
			fmt.Println("NackRequeue")
			return
		case gamelogic.WarOutcomeNoUnits:
			at = pubsub.NackDiscard
			fmt.Println("NackDiscard")
			return
		case gamelogic.WarOutcomeYouWon:
			at = pubsub.Ack
			fmt.Println("Ack")
			return
		case gamelogic.WarOutcomeDraw:
			at = pubsub.Ack
			fmt.Println("Ack")
			return
		default:
			fmt.Println("Error, unrecognized War Outcome")
			at = pubsub.NackDiscard
			fmt.Println("NackDiscard")
			return
		}
	}
}
