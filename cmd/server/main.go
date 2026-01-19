package main

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/havokmoobii/peril/internal/gamelogic"
	"github.com/havokmoobii/peril/internal/pubsub"
	"github.com/havokmoobii/peril/internal/routing"
)

func main() {
	rabbitConnString := "amqp://guest:guest@localhost:5672/"
	
	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf("Could not connect to RabbitMQ: %v\n", err)
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("Could not open channel: %v\n", err)
	}

	_, _, err = pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		routing.GameLogSlug + ".*",
		pubsub.Durable,
	)
	if err != nil {
		log.Fatalf("Could not create queue: %v\n", err)
	}

	gamelogic.PrintServerHelp()

	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}
		switch words[0] {
		case "pause":
			fmt.Println("Sending a pause message.")
			err = pubsub.PublishJSON(channel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
				IsPaused: true,
			})
			if err != nil {
				log.Fatalf("Could not publish time: %v\n", err)
			}
			fmt.Println("Pause message sent!")
		case "resume":
			fmt.Println("Sending a resume message.")
			err = pubsub.PublishJSON(channel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{})
			if err != nil {
				log.Fatalf("Could not publish time: %v\n", err)
			}
			fmt.Println("Resume message sent!")
		case "quit":
			return
		case "help":
			gamelogic.PrintServerHelp()
		default:
			fmt.Println("Unknown command")
		}
	}
}
