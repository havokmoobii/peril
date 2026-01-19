package main

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/havokmoobii/peril/internal/gamelogic"
	"github.com/havokmoobii/peril/internal/routing"
	"github.com/havokmoobii/peril/internal/pubsub"
)

func main() {
	rabbitConnString := "amqp://guest:guest@localhost:5672/"

	fmt.Println("Connecting to RabbitMQ...")
	
	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf("Could not connect to RabbitMQ: %v\n", err)
	}
	defer conn.Close()

	fmt.Println("Peril game client connected to RabbitMQ!")

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("Could not get username: %v\n", err)
	}

	fmt.Println("Username:", username)

	queueName := fmt.Sprintf("%v.%v", routing.PauseKey, username)

	_, _, err = pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilDirect,
		queueName,
		routing.PauseKey,
		pubsub.Transient,
	)
	if err != nil {
		log.Fatalf("Could not create queue: %v\n", err)
	}

	gamestate := gamelogic.NewGameState(username) 

	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}
		switch words[0] {
		case "spawn":
			err = gamestate.CommandSpawn(words)
			if err != nil {
				fmt.Println(err)
			}
		case "move":
			_, err = gamestate.CommandMove(words)
			if err != nil {
				fmt.Println(err)
			}
		case "status":
			gamestate.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			fmt.Println("Spamming not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			return
		default:
			fmt.Println("Unknown command")
		}
	}
}
