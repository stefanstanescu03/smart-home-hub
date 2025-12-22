package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"time"
)

func connect(dialer net.Dialer) net.Conn {
	var conn net.Conn
	var err error

	conn, err = dialer.Dial("tcp", "127.0.0.1:5001")
	if err != nil {
		fmt.Println("[IOT-CLIENT-debug] Connection failed:", err)
	}

	return conn
}

func main() {

	ident := flag.String("id", "", "Device identifier (required)")
	flag.Parse()

	if *ident == "" {
		fmt.Println("ERROR: device id is required")
		fmt.Println("Usage: ./iot-client -id sensor-000381")
		return
	}

	localAddr := &net.TCPAddr{
		IP: net.ParseIP("0.0.0.0"),
	}

	dialer := net.Dialer{
		LocalAddr: localAddr,
	}

	conn := connect(dialer)

	fmt.Println("[IOT-CLIENT-debug] IOT-CLIENT started on port: ", localAddr.Port)

	for {

		temp := 20 + rand.Float32()*10
		humid := 40 + rand.Float32()*30
		message := fmt.Sprintf("ident:%s,temperature[C]:%.2f,humidity[%%]:%.2f\n", *ident, temp, humid)

		if conn != nil {
			_, err := conn.Write([]byte(message))
			if err != nil {
				fmt.Println("Error sending:", err)
				conn = connect(dialer)
			}
		} else {
			fmt.Println("No connection")
			conn = connect(dialer)
		}

		time.Sleep(5 * time.Second)
	}

}
