package utils

import (
	"fmt"
	"os"
	"strings"
)

func ParseMessage(msg string, filename string) {

	// First we remove new line character from our message
	// We split it in multiple parts, separated by ,
	msg = strings.TrimSuffix(msg, "\n")
	parts := strings.Split(msg, ",")

	// We form the csv first row and values that needs to be added to the file
	csv_types := ""
	csv_values := ""

	for _, part := range parts {
		splits := strings.Split(part, ":")
		csv_types = csv_types + splits[0] + ","
		csv_values = csv_values + splits[1] + ","
	}

	csv_types += "\n"
	csv_values += "\n"

	//

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	info, _ := os.Stat(filename)

	if info.Size() == 0 {
		if _, err := f.Write([]byte(csv_types)); err != nil {
			fmt.Println("Error: ", err)
			return
		}
	}

	if _, err := f.Write([]byte(csv_values)); err != nil {
		fmt.Println("Error: ", err)
		return
	}

	if err := f.Close(); err != nil {
		fmt.Println("Error: ", err)
		return
	}
}
