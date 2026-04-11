package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

var httpClient = &http.Client{
	Timeout: time.Second * 10,
}

func SendToCloud(msg string, cloud_api string) error {
	var records []SenMLRecord
	err := json.Unmarshal([]byte(msg), &records)
	if err != nil {
		return fmt.Errorf("failed to parse SenML JSON: %w", err)
	}
	url := cloud_api
	for _, rec := range records {
		placeholder := "{" + rec.Name + "}"
		valStr := fmt.Sprintf("%.2f", rec.Value)
		url = strings.ReplaceAll(url, placeholder, valStr)
	}

	resp, err := httpClient.Get(url)
	if err != nil {
		return fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("cloud API returned status: %d", resp.StatusCode)
	}

	return nil
}
