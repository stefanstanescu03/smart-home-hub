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

	state1 := false
	state2 := false

	var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		if string(msg.Payload()) == "ON_RELAY1" {
			state1 = true
		} else if string(msg.Payload()) == "OFF_RELAY1" {
			state1 = false
		} else if string(msg.Payload()) == "ON_RELAY2" {
			state2 = true
		} else if string(msg.Payload()) == "OFF_RELAY2" {
			state2 = false
		}

		topic := "stat/" + clientID + "/state1"
		token := client.Publish(topic, 1, false, strconv.FormatBool(state1))
		token.Wait()

		topic = "stat/" + clientID + "/state2"
		token = client.Publish(topic, 1, false, strconv.FormatBool(state2))
		token.Wait()

		fmt.Printf("Current state 1: %t\n", state1)
		fmt.Printf("Current state 2: %t\n", state2)
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
