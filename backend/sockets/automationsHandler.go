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

type automationMetadata struct {
	automation    models.Automation
	csv_location  string
	ident_device2 string
}

type routineMetadata struct {
	routine      models.Routine
	ident_device string
}

var automationsPool []automationMetadata
var routinesPool []routineMetadata

var refreshAutomations = make(chan struct{}, 1)

func NotifyAutomationsHandler() {
	select {
	case refreshAutomations <- struct{}{}:
	default:
	}
}

func fetchCurrentAutomations() {

	automationsPool = nil

	var automations []models.Automation
	initializers.DB.Find(&automations)

	for _, automation := range automations {

		var device1 models.Device
		initializers.DB.First(&device1, "id = ?", automation.Device1Id)

		var device2 models.Device
		initializers.DB.First(&device2, "id = ?", automation.Device2Id)

		automationsPool = append(automationsPool, automationMetadata{
			automation:    automation,
			csv_location:  device1.Csv_location,
			ident_device2: device2.Ident,
		})
	}

}

func fetchCurrentRoutines() {

	var routines []models.Routine
	initializers.DB.Find(&routines)

	for _, routine := range routines {

		var device models.Device
		initializers.DB.First(&device, "id = ?", routine.DeviceId)

		routinesPool = append(routinesPool, routineMetadata{
			routine:      routine,
			ident_device: device.Ident,
		})
	}

}

func shouldTriggerRoutine(routine routineMetadata) bool {

	now := time.Now()
	currentMinutes := now.Hour()*60 + now.Minute()
	var startH, startM int
	fmt.Sscanf(routine.routine.StartTime, "%d:%d", &startH, &startM)
	startMinutes := startH*60 + startM

	return currentMinutes >= startMinutes

}

func showldTriggerAutomation(automation automationMetadata) bool {

	// Check for schedule
	start_time := automation.automation.ScheduleStart
	end_time := automation.automation.ScheduleEnd

	if len(start_time) != 0 && len(end_time) != 0 {

		now := time.Now()
		currentMinutes := now.Hour()*60 + now.Minute()

		var startH, startM int
		fmt.Sscanf(start_time, "%d:%d", &startH, &startM)
		startMinutes := startH*60 + startM

		var endH, endM int
		fmt.Sscanf(end_time, "%d:%d", &endH, &endM)
		endMinutes := endH*60 + endM

		if startMinutes <= endMinutes {
			if !(currentMinutes >= startMinutes && currentMinutes <= endMinutes) {
				return false
			}
		} else {
			if !(currentMinutes >= startMinutes || currentMinutes <= endMinutes) {
				return false
			}
		}
	}

	//

	conditionKey, op, conditionValue, ok := splitCondition(automation.automation.Condition)
	if !ok {
		utils.WriteToLogs("[AUTOMATIONS-HANDLER]", fmt.Sprintf("error parsing automation with id %d condition", automation.automation.ID))
		return false
	}

	entry := utils.FetchFromCSV(automation.csv_location)

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
			cleanValue := strings.TrimSpace(value)
			valueFloat, err := strconv.ParseFloat(cleanValue, 64)
			if err != nil {
				utils.WriteToLogs("[AUTOMATIONS-HANDLER]", fmt.Sprintf("error checking automation with id %d: %s", automation.automation.ID, err))
				return false
			}

			conditionValueFloat, err := strconv.ParseFloat(conditionValue, 64)
			if err != nil {
				utils.WriteToLogs("[AUTOMATIONS-HANDLER]", fmt.Sprintf("error checking automation with id %d: %s", automation.automation.ID, err))
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

func StartAutomationsHandler() {

	utils.WriteToLogs("[AUTOMATIONS-HANDLER]", "Started automations handler")

	fetchCurrentAutomations()
	fetchCurrentRoutines()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	automationsStates := make(map[uint]bool)
	routineStates := make(map[uint]bool)

	for {
		select {
		case <-ticker.C:

			for _, routine := range routinesPool {
				if shouldTriggerRoutine(routine) {
					if !routineStates[routine.routine.ID] {

						PublishCmd(routine.ident_device, routine.routine.Payload)

						routineStates[routine.routine.ID] = true
					}
				} else {
					routineStates[routine.routine.ID] = false
				}
			}

			for _, automation := range automationsPool {
				if showldTriggerAutomation(automation) {
					if !automationsStates[automation.automation.ID] {

						PublishCmd(automation.ident_device2, automation.automation.Payload)

						automationsStates[automation.automation.ID] = true
					}
				} else {
					automationsStates[automation.automation.ID] = false
				}
			}
		case <-refreshAutomations:
			fetchCurrentAutomations()
			fetchCurrentRoutines()
			utils.WriteToLogs("[AUTOMATIONS-HANDLER]", "Automations cache refreshed")
		}
	}

}
