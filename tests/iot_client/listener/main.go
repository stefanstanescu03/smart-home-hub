package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {

	ident := flag.String("id", "", "Device identifier (required)")
	flag.Parse()

	if *ident == "" {
		fmt.Println("ERROR: device id is required")
		return
	}

	broker := "tcp://0.0.0.0:5001"
	clientID := *ident
	topic := "cmd/" + clientID

	opts := mqtt.NewClientOptions().AddBroker(broker).SetClientID(clientID)
	opts.SetAutoReconnect(true)

	state := false

	var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		if string(msg.Payload()) == "ON" {
			state = true
		} else if string(msg.Payload()) == "OFF" {
			state = false
		}

		topic := "stat/" + clientID
		token := client.Publish(topic, 1, false, strconv.FormatBool(state))
		token.Wait()

		fmt.Printf("Current state: %t\n", state)
	}

	opts.SetDefaultPublishHandler(messageHandler)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := client.Subscribe(topic, 1, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	fmt.Printf("Listening on topic: %s\n", topic)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\nDisconnecting...")
	client.Disconnect(250)
}
