package sockets

import (
	"backend/initializers"
	"backend/models"
	"backend/utils"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type alertHandlerMetadata struct {
	alert        models.Alert
	csv_location string
}

var alertsPool []alertHandlerMetadata
var refreshAlerts = make(chan struct{}, 1)

func NotifyAlertsHandler() {
	select {
	case refreshAlerts <- struct{}{}:
	default:
	}
}

func fetchCurrentAlerts() {

	alertsPool = nil

	var alerts []models.Alert
	initializers.DB.Find(&alerts)

	for _, alert := range alerts {

		var device models.Device

		initializers.DB.First(&device, "id = ?", alert.DeviceId)

		alertsPool = append(alertsPool, alertHandlerMetadata{
			alert:        alert,
			csv_location: device.Csv_location,
		})
	}
}

func splitCondition(s string) (field, op, value string, ok bool) {
	operators := []string{"<", ">", "="}

	for _, operator := range operators {
		if i := strings.Index(s, operator); i != -1 {
			return strings.TrimSpace(s[:i]),
				operator,
				strings.TrimSpace(s[i+len(operator):]),
				true
		}
	}

	return "", "", "", false
}

func handleSendEmail(alert models.Alert) {

	// Check if you should send it
	var user models.User
	initializers.DB.First(&user, "id = ?", alert.UserId)
	if user.ID == 0 {
		utils.WriteToLogs("ALERTS-HANDLER", fmt.Sprintf("error sending email for alert: %d, can't find the desired user", alert.ID))
		return
	}

	email := user.Email

	if alert.NotifyEmail {
		ok := utils.SendSimpleEmail(email, alert.Subject, alert.Content)
		if !ok {
			utils.WriteToLogs("ALERTS-HANDLER", fmt.Sprintf("error sending email for alert: %d", alert.ID))
		}
	}

}

func shouldTrigger(alertMetadata alertHandlerMetadata) bool {

	conditionKey, op, conditionValue, ok := splitCondition(alertMetadata.alert.Condition)
	if !ok {
		utils.WriteToLogs("ALERTS-HANDLER", fmt.Sprintf("error parsing alert with id %d condition", alertMetadata.alert.ID))
		return false
	}

	entry := utils.FetchFromCSV(alertMetadata.csv_location)

	if entry == "-1" {
		return false
	}

	parts := strings.Split(entry, ",")

	for _, part := range parts {

		kv := strings.SplitN(part, ":", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.SplitN(kv[0], "[", 2)[0]
		value := kv[1]

		if key == conditionKey {
			valueFloat, err := strconv.ParseFloat(value, 64)
			if err != nil {
				utils.WriteToLogs("ALERTS-HANDLER", fmt.Sprintf("error checking alert with id %d: %s", alertMetadata.alert.ID, err))
				return false
			}

			conditionValueFloat, err := strconv.ParseFloat(conditionValue, 64)
			if err != nil {
				utils.WriteToLogs("ALERTS-HANDLER", fmt.Sprintf("error checking alert with id %d: %s", alertMetadata.alert.ID, err))
				return false
			}

			switch op {
			case "<":
				return valueFloat < conditionValueFloat
			case ">":
				return valueFloat > conditionValueFloat
			case "=":
				return valueFloat == conditionValueFloat
			default:
				return false
			}
		}
	}
	return false
}

func StartAlertsHandler() {

	utils.WriteToLogs("[ALERTS-HANDLER]", "Started alerts handler")
	fetchCurrentAlerts()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	alertStates := make(map[uint]bool)

	for {
		select {
		case <-ticker.C:
			for _, alert := range alertsPool {
				if shouldTrigger(alert) {
					if !alertStates[alert.alert.ID] {
						utils.WriteAlert(alert.alert.DeviceId, alert.alert.Subject, alert.alert.Content)
						handleSendEmail(alert.alert)
						alertStates[alert.alert.ID] = true
					}
				} else {
					alertStates[alert.alert.ID] = false
				}
			}
		case <-refreshAlerts:
			fetchCurrentAlerts()
			utils.WriteToLogs("[ALERTS-HANDLER]", "Alerts cache refreshed")
		}
	}
}
