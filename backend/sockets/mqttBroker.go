package sockets

import (
	"backend/initializers"
	"backend/models"
	"backend/utils"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	mqtt "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/listeners"
	"github.com/mochi-mqtt/server/v2/packets"
)

type TelemetryHook struct {
	mqtt.HookBase
}

type ConnectionMetadata struct {
	Csv_location string
	Cloud_api    string
}

var ConnectionPool sync.Map

func (h *TelemetryHook) ID() string {
	return "telemetry-csv-logger"
}

func (h *TelemetryHook) Provides(b byte) bool {
	return b == mqtt.OnConnectAuthenticate || b == mqtt.OnPublish || b == mqtt.OnDisconnect || b == mqtt.OnACLCheck
}

func (h *TelemetryHook) OnConnectAuthenticate(cl *mqtt.Client, pk packets.Packet) bool {

	if cl.ID != "SYSTEM_CMD_CLIENT" {
		ident := string(cl.ID)

		var device models.Device
		initializers.DB.First(&device, "ident = ?", ident)
		if device.ID == 0 {
			utils.WriteToLogs("[MQTT-BROKER]", fmt.Sprintf("Connection rejected: %s", ident))
			return false
		}

		ConnectionPool.Store(cl.ID, ConnectionMetadata{Csv_location: device.Csv_location, Cloud_api: device.Cloud_api})

		utils.WriteToLogs("[MQTT-BROKER]", fmt.Sprintf("New device connected: %s", ident))
	}
	return true
}

func (h *TelemetryHook) OnPublish(cl *mqtt.Client, pk packets.Packet) (packets.Packet, error) {

	topic := pk.TopicName

	if metadata, ok := ConnectionPool.Load(cl.ID); ok && strings.HasPrefix(topic, "telemetry/") {
		msg := string(pk.Payload)
		utils.ParseMessage(msg, metadata.(ConnectionMetadata).Csv_location)

		if len(metadata.(ConnectionMetadata).Cloud_api) != 0 {
			err := utils.SendToCloud(msg, metadata.(ConnectionMetadata).Cloud_api)
			if err != nil {
				utils.WriteToLogs("[MQTT-BROKER]", fmt.Sprintf("Error sending %s data to cloud: %s", cl.ID, err))
			}
		}

	}
	return pk, nil
}

func (h *TelemetryHook) OnACLCheck(cl *mqtt.Client, topic string, write bool) bool {

	if cl.ID == "SYSTEM_CMD_CLIENT" {
		return true
	}

	if write {
		if strings.HasPrefix(topic, fmt.Sprintf("telemetry/%s", cl.ID)) || strings.HasPrefix(topic, fmt.Sprintf("stat/%s", cl.ID)) {
			return true
		}

		utils.WriteToLogs("[MQTT-BROKER]", fmt.Sprintf("A client tried to publish in a forbidden topic: %s", cl.ID))
		return false
	}

	if strings.HasPrefix(topic, fmt.Sprintf("cmd/%s", cl.ID)) {
		return true
	}

	utils.WriteToLogs("[MQTT-BROKER]", fmt.Sprintf("A client tried to subscribe to a forbidden topic: %s", cl.ID))
	return false

}

func (h *TelemetryHook) OnDisconnect(cl *mqtt.Client, err error, reason bool) {

	ConnectionPool.Delete(cl.ID)
	utils.WriteToLogs("[MQTT-BROKER]", fmt.Sprintf("Device disconnected: %s", cl.ID))
}

func StartMQTTBroker() {

	server := mqtt.New(nil)

	hook := &TelemetryHook{}

	err := server.AddHook(hook, nil)
	if err != nil {
		utils.WriteToLogs("[MQTT-BROKER]", "Error starting MQTT Broker")
		log.Fatal(err)
	}

	// _ = server.AddHook(new(auth.AllowHook), nil)

	port := os.Getenv("TELEMETRY_PORT")
	host := os.Getenv("HOST")

	tcp := listeners.NewTCP(listeners.Config{ID: "t1", Address: host + ":" + port})
	err = server.AddListener(tcp)
	if err != nil {
		utils.WriteToLogs("[MQTT-BROKER]", "Error starting MQTT Broker")
		log.Fatal(err)
	}

	utils.WriteToLogs("[MQTT-BROKER]", fmt.Sprintf("MQTT Broker started on port: %s", port))

	go func() {
		err := server.Serve()
		if err != nil {
			log.Fatal(err)
		}
	}()
}
