package sockets

import (
	"backend/utils"
	"fmt"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var cmd_client mqtt.Client

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

func StartCmdClient() {

	port := os.Getenv("TELEMETRY_PORT")
	host := os.Getenv("HOST")

	broker_addr := "tcp://" + host + ":" + port

	opts := mqtt.NewClientOptions().AddBroker(broker_addr).SetClientID("SYSTEM_CMD_CLIENT")
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetConnectTimeout(5 * time.Second)

	cmd_client = mqtt.NewClient(opts)

	if token := cmd_client.Connect(); token.Wait() && token.Error() != nil {
		utils.WriteToLogs("[CMD-CLIENT]", "Failed to connect to broker")
		panic(token.Error())
	}

	utils.WriteToLogs("[CMD-CLIENT]", "Client connected to broker")

}
