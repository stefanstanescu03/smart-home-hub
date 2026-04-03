package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type SenMLRecord struct {
	Name  string  `json:"n"`
	Unit  string  `json:"u"`
	Value float64 `json:"v"`
}

func ParseMessage(msg string, filename string) {

	var records []SenMLRecord
	err := json.Unmarshal([]byte(msg), &records)
	if err != nil {
		fmt.Println("Error: ", err)
	}

	var csvTypes strings.Builder
	var csvData strings.Builder

	currTime := time.Now().Unix()
	csvTypes.WriteString("timestamp,")
	csvData.WriteString(strconv.FormatInt(currTime, 10) + ",")

	for _, record := range records {
		header := fmt.Sprintf("%s[%s]", record.Name, record.Unit)
		csvTypes.WriteString(header + ",")
		csvData.WriteString(strconv.FormatFloat(record.Value, 'f', 2, 64) + ",")
	}

	typesStr := strings.TrimSuffix(csvTypes.String(), ",") + "\n"
	valuesStr := strings.TrimSuffix(csvData.String(), ",") + "\n"

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	info, _ := os.Stat(filename)

	if info.Size() == 0 {
		if _, err := f.Write([]byte(typesStr)); err != nil {
			fmt.Println("Error: ", err)
			return
		}
	}

	if _, err := f.Write([]byte(valuesStr)); err != nil {
		fmt.Println("Error: ", err)
		return
	}

	if err := f.Close(); err != nil {
		fmt.Println("Error: ", err)
		return
	}
}
