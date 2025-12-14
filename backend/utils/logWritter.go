package utils

import (
	"fmt"
	"os"
	"time"
)

func WriteToLogs(sender string, message string) {

	logsLocation := os.Getenv("LOGS") + "/sys.log"
	currTime := time.Now().Format(time.RFC3339)

	log := currTime + " " + sender + " " + message + "\n"
	fmt.Print(log)

	f, err := os.OpenFile(logsLocation, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	if _, err := f.Write([]byte(log)); err != nil {
		fmt.Println("Error writing to logs file: ", err)
		return
	}

	if err := f.Close(); err != nil {
		fmt.Println("Error closing logs file: ", err)
		return
	}
}

func WriteAlert(deviceId uint, subject string, content string) {

	logsLocation := os.Getenv("LOGS") + "/alerts.log"
	currTime := time.Now().Format(time.RFC3339)

	log := currTime + fmt.Sprintf(" [ALERT:%d] ", deviceId) + subject + " : " + content + "\n"

	f, err := os.OpenFile(logsLocation, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	if _, err := f.Write([]byte(log)); err != nil {
		fmt.Println("Error writing to alerts file: ", err)
		return
	}

	if err := f.Close(); err != nil {
		fmt.Println("Error closing alerts file: ", err)
		return
	}
}
