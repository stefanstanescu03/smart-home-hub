package sockets

import (
	"backend/initializers"
	"backend/models"
	"backend/utils"
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

var TelemetryConnectionPool sync.Map

func handleConnection(conn net.Conn) {

	defer conn.Close()

	reader := bufio.NewReader(conn)

	// Check if the device is registered
	firstMessage, err := reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			return
		}
		utils.WriteToLogs("[TELEMETRY]", fmt.Sprintf("Error reading a message from %s: %s", conn.RemoteAddr(), err))
		return
	}
	kv := strings.Split(strings.Split(firstMessage, ",")[0], ":")
	template := kv[0]
	ident := kv[1]

	if template != "ident" {
		utils.WriteToLogs("[TELEMETRY]", fmt.Sprintf("Failed to retrieve the identificator from the device from ip %s: %s", conn.RemoteAddr(), err))
		return
	}

	if _, exists := TelemetryConnectionPool.LoadOrStore(ident, true); exists {
		utils.WriteToLogs("[TELEMETRY]", fmt.Sprintf("Devices with the same identificator: %s", ident))
		return
	}

	var device models.Device
	initializers.DB.First(&device, "ident = ?", ident)
	if device.ID == 0 {
		utils.WriteToLogs("[TELEMETRY]", fmt.Sprintf("Connection rejected: %s", ident))
		return
	}

	//

	defer func() {
		utils.WriteToLogs("[TELEMETRY]", fmt.Sprintf("Device disconnected: %s", ident))
		TelemetryConnectionPool.Delete(ident)
	}()

	utils.WriteToLogs("[TELEMETRY]", fmt.Sprintf("New device connected: %s", ident))
	TelemetryConnectionPool.Store(ident, true)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return
			}
			utils.WriteToLogs("[TELEMETRY]", fmt.Sprintf("Error reading a message from %s: %s", conn.RemoteAddr(), err))
			return
		}
		utils.ParseMessage(msg, device.Csv_location)
	}

}

func StartTelemetryServer() {

	port := os.Getenv("TELEMETRY_PORT")
	host := os.Getenv("HOST")
	listener, err := net.Listen("tcp", host+":"+port)

	if err != nil {
		log.Fatal(err)
	}

	utils.WriteToLogs("[TELEMETRY]", fmt.Sprintf("Telemetry server started on port: %s", port))

	for {
		conn, err := listener.Accept()
		if err != nil {
			utils.WriteToLogs("[TELEMETRY]", fmt.Sprintf("Error accepting connection: %s", err))
			continue
		}

		go handleConnection(conn)

	}
}
