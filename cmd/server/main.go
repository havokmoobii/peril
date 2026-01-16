package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	amqp "github.com/rabbitmq/amqp091-go"

	"github.com/havokmoobii/peril/internal/pubsub"
	"github.com/havokmoobii/peril/internal/routing"
)

func main() {
	rabbitConnString := "amqp://guest:guest@localhost:5672/"

	fmt.Println("Connecting to RabbitMQ...")
	
	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf("Could not connect to RabbitMQ: %v\n", err)
	}
	defer conn.Close()

	fmt.Println("Peril game server connected to RabbitMQ!")

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("Could not open channel: %v\n", err)
	}

	err = pubsub.PublishJSON(channel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{
		IsPaused: true,
	})
	if err != nil {
		log.Fatalf("Could not publish time: %v\n", err)
	}
	fmt.Println("Pause message sent!")

	// Wait for ctrl+c
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan

	fmt.Println("\nRabbitMQ connection closed.")
}
