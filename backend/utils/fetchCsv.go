package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func FetchFromCSV(filename string) string {
	f, err := os.Open(filename)

	if err != nil {
		fmt.Println("Error: ", err)
		return ""
	}

	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println("Error: ", err)
		}
	}()

	scanner := bufio.NewScanner(f)
	if scanner.Scan() {
		csv_types := scanner.Text()
		csv_types = strings.Trim(csv_types, ",") // entries will usually have an useless , at the end
		types := strings.Split(csv_types, ",")

		if err := scanner.Err(); err != nil {
			fmt.Println("Error scanning file:", err)
			return ""
		}

		// Read last line
		line := ""
		// var cursor int64 = 0
		// stat, _ := f.Stat()
		// filesize := stat.Size()
		lineBytes := make([]byte, 0, 100) // preallocate a slice for efficiency
		buf := make([]byte, 1)            // reuse buffer for reading

		for cursor := int64(-1); ; cursor-- {
			_, err := f.Seek(cursor, io.SeekEnd)
			if err != nil {
				fmt.Println("Seek error:", err)
				break
			}

			_, err = f.Read(buf)
			if err != nil {
				if err == io.EOF {
					break
				}
				fmt.Println("Read error:", err)
				break
			}

			// Stop if we find a newline (but skip trailing newline at end)
			if cursor != -1 && (buf[0] == '\n' || buf[0] == '\r') {
				break
			}

			// Prepend the byte to the slice
			lineBytes = append([]byte{buf[0]}, lineBytes...)

			// Stop if we reach the beginning of the file
			stat, _ := f.Stat()
			if -cursor >= stat.Size() {
				break
			}
		}

		// Convert bytes to string
		line = string(lineBytes)
		line = strings.Trim(line, ",") // entries will usually have an useless , at the end
		values := strings.Split(line, ",")
		final_str := ""
		for i := range types {
			final_str = final_str + types[i] + ":" + values[i] + ","
		}
		final_str = strings.Trim(final_str, ",")
		return final_str

	} else {
		fmt.Println("Empty file")
		return ""
	}
}
