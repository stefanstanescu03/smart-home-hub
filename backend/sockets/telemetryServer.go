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
	"time"
)

func handleConnection(conn net.Conn, csv_location string) {

	defer func() {
		fmt.Println("Device disconnected: ", conn.RemoteAddr())
		conn.Close()
	}()

	reader := bufio.NewReader(conn)

	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return
			}
			fmt.Println("Read error:", err)
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

	fmt.Println("[TELEMETRY-debug] Telemetry server started on port: ", port)

	for {

		if time.Since(lastTime).Minutes() > 1 {
			lastTime = time.Now()
			devices_cache = make(map[string]bool)
		}

		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err)
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
				fmt.Println("Rejected: ", conn.RemoteAddr())
				conn.Close()
				continue
			}
		}

		fmt.Println("New device connected: ", conn.RemoteAddr())

		go handleConnection(conn, device.Csv_location)

	}
}
