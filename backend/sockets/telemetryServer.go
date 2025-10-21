package sockets

import (
	"backend/initializers"
	"backend/models"
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
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

		f, err := os.OpenFile(csv_location, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}
		if _, err := f.Write([]byte(msg)); err != nil {
			fmt.Println("Error: ", err)
			continue
		}
		if err := f.Close(); err != nil {
			fmt.Println("Error: ", err)
			continue
		}
		// conn.Write([]byte("OK\n"))
	}

}

func StartTelemetryServer() {

	port := os.Getenv("TELEMETRY_PORT")
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("[TELEMETRY-debug] Telemetry server started on port: ", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err)
			continue
		}

		var device models.Device
		initializers.DB.First(&device, "ip = ?", conn.RemoteAddr().String())

		if device.ID == 0 {
			fmt.Println("Rejected: ", conn.RemoteAddr())
			conn.Close()
			continue
		}

		fmt.Println("New device connected: ", conn.RemoteAddr())

		go handleConnection(conn, device.Csv_location)

	}
}
