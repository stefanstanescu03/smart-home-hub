package sockets

import (
	"backend/utils"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type DeviceKey struct {
	Ident   string
	Channel string
}

var cmd_client mqtt.Client
var DeviceStates sync.Map

func PublishCmd(ident string, payload string) {
	if cmd_client == nil || !cmd_client.IsConnected() {
		utils.WriteToLogs("[CMD-CLIENT]", fmt.Sprintf("Failed to publish command: %s", payload))
	}
	topic := fmt.Sprintf("cmd/%s", ident)
	token := cmd_client.Publish(topic, 1, false, payload)

	go func() {
		_ = token.Wait()
		if token.Error() != nil {
			utils.WriteToLogs("[CMD-CLIENT]", fmt.Sprintf("Failed to publish command: %s, error: %s", payload, token.Error()))
		}
	}()
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	topic := msg.Topic()
	payload := string(msg.Payload())

	if strings.HasPrefix(topic, "stat/") {

		parts := strings.Split(topic, "/")
		ident := parts[1]

		var channel string

		if len(parts) > 2 {
			channel = parts[2]
		} else {
			channel = "unregistered"
		}

		DeviceStates.Store(DeviceKey{Ident: ident, Channel: channel}, payload)
	}
}

func StartCmdClient() {

	port := os.Getenv("TELEMETRY_PORT")
	host := os.Getenv("HOST")

	broker_addr := "tcp://" + host + ":" + port

	opts := mqtt.NewClientOptions().AddBroker(broker_addr).SetClientID("SYSTEM_CMD_CLIENT")
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetConnectTimeout(5 * time.Second)
	opts.SetDefaultPublishHandler(messagePubHandler)

	cmd_client = mqtt.NewClient(opts)

	if token := cmd_client.Connect(); token.Wait() && token.Error() != nil {
		utils.WriteToLogs("[CMD-CLIENT]", "Failed to connect to broker")
		panic(token.Error())
	}

	if token := cmd_client.Subscribe("stat/#", 1, nil); token.Wait() && token.Error() != nil {
		utils.WriteToLogs("[CMD-CLIENT]", "Client failed to subscribe to stat/#")
		panic(token.Error())
	}

	utils.WriteToLogs("[CMD-CLIENT]", "Client connected to broker")

}
