package main

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

func connect(dialer net.Dialer) net.Conn {
	var conn net.Conn
	var err error

	conn, err = dialer.Dial("tcp", "host.docker.internal:5001")
	if err != nil {
		fmt.Println("[IOT-CLIENT-debug] Connection failed:", err)
	}

	return conn
}

func main() {

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
		iaq := rand.Float32() * 200
		dust := 10 + rand.Float32()*90
		message := fmt.Sprintf("temperature[C]:%.2f,humidity[%%]:%.2f,IAQ:%.2f,Dust[ug/m^3]:%.2f\n", temp, humid, iaq, dust)

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
