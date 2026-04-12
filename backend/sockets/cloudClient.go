package sockets

import (
	"backend/utils"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type BlynkBridge struct {
	clients sync.Map
}

var Bridge = &BlynkBridge{}

func (b *BlynkBridge) GetClient(token string) mqtt.Client {
	if val, ok := b.clients.Load(token); ok {
		client := val.(mqtt.Client)
		if client.IsConnected() {
			return client
		}
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker("tls://blynk.cloud:8883")
	opts.SetUsername("device")
	opts.SetPassword(token)
	opts.SetClientID("Hub_Bridge_" + token)

	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(1 * time.Minute)

	client := mqtt.NewClient(opts)
	if t := client.Connect(); t.Wait() && t.Error() != nil {
		utils.WriteToLogs("[BLYNK]", fmt.Sprintf("Connection failed for %s: %v", token, t.Error()))
		return nil
	}

	utils.WriteToLogs("[BLYNK]", fmt.Sprintf("Connected successfully for token: %s", token))
	b.clients.Store(token, client)
	return client
}

func (b *BlynkBridge) PushData(token string, pin int, value string) {
	client := b.GetClient(token)
	if client == nil {
		return
	}
	var measurements []utils.SenMLRecord
	err := json.Unmarshal([]byte(value), &measurements)
	if err != nil {
		log.Printf("[BLYNK] Error parsing SenML: %v", err)
		return
	}

	for i := range measurements {
		topic := fmt.Sprintf("ds/%s", measurements[i].Name)
		payload := fmt.Sprintf("%.2f", measurements[i].Value)
		tokenPublish := client.Publish(topic, 1, false, payload)

		go func() {
			_ = tokenPublish.Wait()
			if tokenPublish.Error() != nil {
				utils.WriteToLogs("[BLYNK]", fmt.Sprintf("Failed to publish: %s, error: %s", payload, tokenPublish.Error()))
			}
		}()
	}
}
