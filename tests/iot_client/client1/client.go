package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"net"
	"time"
)

func main() {

	localAddr := &net.TCPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: 40000,
	}

	dialer := net.Dialer{
		LocalAddr: localAddr,
	}

	conn, err := dialer.Dial("tcp", "localhost:5001")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("[IOT-CLIENT-debug] IOT-CLIENT started on port: ", localAddr.Port)

	for {

		temp := 20 + rand.Float32()*10
		humid := 40 + rand.Float32()*30
		iaq := rand.Float32() * 200
		dust := 10 + rand.Float32()*90
		message := fmt.Sprintf("temperature[C]:%.2f,humidity[%%]:%.2f,IAQ:%.2f,Dust[ug/m^3]:%.2f\n", temp, humid, iaq, dust)

		_, err := conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error sending:", err)
			break
		}

		time.Sleep(5 * time.Second)
	}

}
