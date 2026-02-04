package utils

import (
	"bufio"
	"bytes"
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

func GetLastLines(filepath string, lineCount int) ([]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stat, _ := file.Stat()
	filesize := stat.Size()

	var lines []string
	var cursor int64 = 0
	bufferSize := int64(1024)
	trailingData := []byte{}

	for len(lines) < lineCount && abs(cursor) < filesize {
		// Calculate how much to read
		currReadSize := bufferSize
		if abs(cursor)+bufferSize > filesize {
			currReadSize = filesize - abs(cursor)
		}

		cursor -= currReadSize
		_, err := file.Seek(cursor, io.SeekEnd)
		if err != nil {
			return nil, err
		}

		buf := make([]byte, currReadSize)
		_, err = io.ReadFull(file, buf)
		if err != nil {
			return nil, err
		}

		// Prepend the new chunk to what we had left over
		combined := append(buf, trailingData...)
		parts := bytes.Split(combined, []byte("\n"))

		// The first part of the split is likely an incomplete line
		trailingData = parts[0]

		// Add parts in reverse order (excluding the incomplete first part)
		for i := len(parts) - 1; i >= 1; i-- {
			if len(parts[i]) == 0 && i == len(parts)-1 {
				continue // Skip trailing newline at the very end of file
			}
			lines = append(lines, string(parts[i]))
			if len(lines) == lineCount {
				break
			}
		}
	}

	// Add the final remaining bit if we haven't hit the limit
	if len(lines) < lineCount && len(trailingData) > 0 {
		lines = append(lines, string(trailingData))
	}

	// REVERSE: Flip the slice back to original file order
	for i, j := 0, len(lines)-1; i < j; i, j = i+1, j-1 {
		lines[i], lines[j] = lines[j], lines[i]
	}

	return lines, nil
}

func abs(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}

func GetFirstLine(filepath string) string {

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

	types := GetFirstLine(filename)
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
