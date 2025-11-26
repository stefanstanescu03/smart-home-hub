package sockets

import (
	"backend/initializers"
	"backend/models"
	"backend/utils"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Client struct {
	conn    *websocket.Conn
	devices map[string]string
	mu      sync.Mutex
}

var clients = make(map[*Client]bool)
var clientsMu sync.Mutex

func wsHandler(w http.ResponseWriter, r *http.Request) {

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}

	// Registering a client
	client := &Client{
		conn:    conn,
		devices: make(map[string]string),
	}
	clientsMu.Lock()
	clients[client] = true
	clientsMu.Unlock()
	fmt.Println("[WS-debug] New client connected")

	// Disconnect a client
	defer func() {
		clientsMu.Lock()
		delete(clients, client)
		clientsMu.Unlock()
		conn.Close()
		fmt.Println("[WS-debug] client disconnected")
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}

		// Get the id from the message (that should be like subscribe:id)
		parts := strings.Split(strings.Trim(string(message), "\""), ":")
		if len(parts) != 2 || parts[0] != "subscribe" {
			fmt.Println("[WS-debug] Invalid subscribe message")
			break
		}
		id := parts[1]

		// Get device csv_location
		var device models.Device
		initializers.DB.First(&device, "id = ?", id)
		if device.ID == 0 {
			fmt.Println("[WS-debug] Device not found")
			break
		}
		csv_location := device.Csv_location

		// Add device
		client.mu.Lock()
		client.devices[id] = csv_location
		fmt.Println("[WS-debug] Client subscribed to device: ", id)
		client.mu.Unlock()

		if err := conn.WriteMessage(websocket.TextMessage, []byte(csv_location)); err != nil {
			fmt.Println("Error writing message:", err)
			break
		}
	}

}

func broadcastHandler() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		clientsMu.Lock()
		for client := range clients {
			for device := range client.devices {
				csv_location := client.devices[device]

				// Get desired data from location
				message := utils.FetchFromCSV(csv_location)

				// Add an id so i know where this data belongs
				message = message + ",id:" + device

				err := client.conn.WriteMessage(websocket.TextMessage, []byte(message))
				if err != nil {
					client.conn.Close()
					delete(clients, client)
				}
			}
		}
		clientsMu.Unlock()
	}
}

func StartStreamingServer() {
	http.HandleFunc("/streaming", wsHandler)
	go broadcastHandler()
	fmt.Println("[WS-debug] Listening on port: 5003")
	err := http.ListenAndServe(":5003", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
