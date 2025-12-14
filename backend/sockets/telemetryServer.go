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
	"time"
)

var TelemetryConnectionPool sync.Map
var AcknowledgedDevices = []string{}

func handleConnection(conn net.Conn, csv_location string) {

	defer func() {
		utils.WriteToLogs("[TELEMETRY]", fmt.Sprintf("Device disconnected: %s", conn.RemoteAddr()))
		ip := strings.Split(conn.RemoteAddr().String(), ":")[0]
		TelemetryConnectionPool.Delete(ip)
		conn.Close()
	}()

	reader := bufio.NewReader(conn)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return
			}
			utils.WriteToLogs("[TELEMETRY]", fmt.Sprintf("Error reading a message from %s: %s", conn.RemoteAddr(), err))
			return
		}

		utils.ParseMessage(msg, csv_location)
	}

}

func StartTelemetryServer() {

	port := os.Getenv("TELEMETRY_PORT")
	host := os.Getenv("HOST")
	listener, err := net.Listen("tcp", host+":"+port)

	lastTime := time.Now()

	var devices_cache = make(map[string]bool)
	if err != nil {
		log.Fatal(err)
	}

	utils.WriteToLogs("[TELEMETRY]", fmt.Sprintf("Telemetry server started on port: %s", port))

	for {

		if time.Since(lastTime).Minutes() > 1 {
			lastTime = time.Now()
			devices_cache = make(map[string]bool)
			AcknowledgedDevices = []string{}
		}

		conn, err := listener.Accept()
		if err != nil {
			utils.WriteToLogs("[TELEMETRY]", fmt.Sprintf("Error accepting connection: %s", err))
			continue
		}

		ip := strings.Split(conn.RemoteAddr().String(), ":")[0]
		var device models.Device

		if devices_cache[ip] {
			conn.Close()
			continue
		} else {
			initializers.DB.First(&device, "ip = ?", ip)
			if device.ID == 0 {
				devices_cache[ip] = true
				AcknowledgedDevices = append(AcknowledgedDevices, ip)
				utils.WriteToLogs("[TELEMETRY]", fmt.Sprintf("Connection rejected: %s", conn.RemoteAddr()))
				conn.Close()
				continue
			}
		}

		utils.WriteToLogs("[TELEMETRY]", fmt.Sprintf("New device connected: %s", conn.RemoteAddr()))
		TelemetryConnectionPool.Store(ip, true)

		go handleConnection(conn, device.Csv_location)

	}
}
