package utils

import "strings"

func Format(name string) string {
	formatted_name := strings.Replace(name, " ", "_", -1)
	formatted_name = strings.ToLower(formatted_name)
	return formatted_name
}
