package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// Source - https://stackoverflow.com/a
// Posted by Bactisme
// Retrieved 2025-12-31, License - CC BY-SA 4.0

func getLastLine(filepath string) string {

	fileHandle, err := os.Open(filepath)

	if err != nil {
		return ""
	}
	defer fileHandle.Close()

	line := ""
	var cursor int64 = 0
	stat, _ := fileHandle.Stat()
	filesize := stat.Size()
	for {
		cursor -= 1
		fileHandle.Seek(cursor, io.SeekEnd)

		char := make([]byte, 1)
		fileHandle.Read(char)

		if cursor != -1 && (char[0] == 10 || char[0] == 13) { // stop if we find a line
			break
		}

		line = fmt.Sprintf("%s%s", string(char), line) // there is more efficient way

		if cursor == -filesize { // stop if we are at the begining
			break
		}
	}

	return line
}

func getFirstLine(filepath string) string {

	fileHandle, err := os.Open(filepath)

	if err != nil {
		return ""
	}
	defer fileHandle.Close()

	scanner := bufio.NewScanner(fileHandle)
	scanner.Scan()

	return scanner.Text()

}

func FetchFromCSV(filename string) string {

	types := getFirstLine(filename)
	if types == "" {
		return "-1"
	}

	value := getLastLine(filename)
	if value == types || value == "" {
		return "-1"
	}

	types = strings.Trim(types, ",")
	value = strings.Trim(value, ",")
	type_parts := strings.Split(types, ",")
	value_parts := strings.Split(value, ",")

	final_str := ""

	for i := range type_parts {
		final_str += type_parts[i] + ":" + value_parts[i] + ","
	}

	final_str = strings.Trim(final_str, ",")

	return final_str
}
