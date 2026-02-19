package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"

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

	opts := mqtt.NewClientOptions().AddBroker(broker).SetClientID(clientID)
	opts.SetKeepAlive(60 * time.Second)
	opts.SetAutoReconnect(true)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	fmt.Printf("Connected as %s. Starting data stream...\n", clientID)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for range ticker.C {

		temp := 20 + r.Float64()*(30-20)
		humidity := 40 + r.Intn(20)

		payload := fmt.Sprintf("temperature[C]:%.2f,humidity[%%]:%d", temp, humidity)

		topic := "telemetry/sensors"
		token := client.Publish(topic, 1, false, payload)
		token.Wait()

		fmt.Printf("[%s] Sent: %s\n", time.Now().Format("15:04:05"), payload)
	}
}
