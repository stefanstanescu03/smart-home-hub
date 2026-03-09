package utils

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

var httpClient = &http.Client{
	Timeout: time.Second * 10,
}

func SendToCloud(msg string, cloud_api string) error {

	// Replace placeholders in cloud_api with the actual parameter values
	parts := strings.Split(msg, ",")
	url := cloud_api

	for i := range len(parts) {
		kv := strings.Split(parts[i], ":")
		key := "{" + strings.TrimSpace(kv[0]) + "}"
		val := strings.TrimSpace(kv[1])

		url = strings.ReplaceAll(url, key, val)
	}
	//

	resp, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("cloud returned error: %d", resp.StatusCode)
	}

	return nil

}
