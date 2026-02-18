package pipelines

import (
	"backend/initializers"
	"backend/models"
	"backend/pipelines/core"
	"backend/utils"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type ModelInfo struct {
	model_metadata   *models.AnomalyModel
	model_parameters *core.DistanceBasedDetector
}

var modelsPool []*ModelInfo

var refreshModels = make(chan struct{}, 1)

var alertStates = make(map[uint]bool)

func handleSendEmail(model models.AnomalyModel) {

	var user models.User
	initializers.DB.First(&user, "id = ?", model.UserId)
	if user.ID == 0 {
		utils.WriteToLogs("ALERTS-HANDLER", fmt.Sprintf("error sending email for anomaly model: %d, can't find the desired user", model.ID))
		return
	}

	email := user.Email

	ok := utils.SendSimpleEmail(email, "Anomaly detected", fmt.Sprintf("Unusual behaviour with parameter: %s", model.Param))
	if !ok {
		utils.WriteToLogs("ALERTS-HANDLER", fmt.Sprintf("error sending email for anomaly model: %d", model.ID))
	}

}

func NotifyAnomalyPipeline() {
	select {
	case refreshModels <- struct{}{}:
	default:
	}
}

func create_model(data_path string, window_size int, param string) *core.DistanceBasedDetector {

	model := core.NewAnomalyDetectionModel(1, 1e-2, window_size)
	dataset := core.Csv_to_dataset(data_path, window_size, param)
	core.Fit(model, &dataset)

	return model
}

func save_model(filename string, model *core.DistanceBasedDetector) {

	json_value, err := json.Marshal(model)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(filename, json_value, 0644)
	if err != nil {
		log.Fatal(err)
	}

}

func load_model(filename string) *core.DistanceBasedDetector {

	json_value, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	model := core.NewAnomalyDetectionModel(1, 1e-2, 10)
	json.Unmarshal(json_value, model)

	return model

}

// Function to create the initial model and save it to the desired file
// This should be called in controller
func FitAndSave(filename string, data_path string, param string) {

	model := create_model(data_path, 10, param)
	save_model(filename, model)
	utils.WriteToLogs("[PIPELINES]", fmt.Sprintf("New model created at: %s", filename))

}

// See if you have anomalies
// This calls predict function that also will adapt parameters
func feed_data(model *ModelInfo) {

	// Get the last window from csv file
	window_size := len(model.model_parameters.Centroid)

	device_id := model.model_metadata.DeviceId
	var device models.Device
	initializers.DB.First(&device, "id = ?", device_id)

	if device.ID == 0 {
		utils.WriteToLogs("[PIPELINES]", fmt.Sprintf("Somehow the device was not found: %d", device_id))
		return
	}

	data_location := device.Csv_location

	// Find what parameter are you interested in
	param := model.model_metadata.Param
	param_id := -1
	first_line := utils.GetFirstLine(data_location)
	keys := strings.Split(first_line, ",")
	for i := range keys {
		if keys[i] == param {
			param_id = i
			break
		}
	}
	lines, err := utils.GetLastLines(data_location, window_size)
	if err != nil {
		utils.WriteToLogs("[PIPELINES]", fmt.Sprintf("Failed to get last values in: %s", data_location))
		return
	}

	// Creating the actual window
	var window core.TimeseriesWindow
	window.NumFeatures = window_size

	for i := range lines {
		line := lines[i]
		keys := strings.Split(line, ",")

		value, err := strconv.ParseFloat(keys[param_id], 32)
		if err != nil {
			log.Fatal(err)
		}
		final_value := float32(value)

		window.Window = append(window.Window, float32(final_value))
	}
	//

	if core.Predict(model.model_parameters, window) {
		if !alertStates[model.model_metadata.ID] {
			utils.WriteAlert(model.model_metadata.DeviceId, "Anomaly detected", fmt.Sprintf("Unusual behaviour with parameter: %s", param))
			if model.model_metadata.NotifyEmail {
				handleSendEmail(*model.model_metadata)
			}
			alertStates[model.model_metadata.ID] = true
		}
	} else {
		alertStates[model.model_metadata.ID] = false
		save_model(model.model_metadata.Location, model.model_parameters)
	}

}

func fetchModels() {

	modelsPool = nil

	var models_metadata []models.AnomalyModel
	initializers.DB.Find(&models_metadata)

	for _, model := range models_metadata {

		data_location := model.Location
		fetched_model_info := ModelInfo{
			model_metadata:   &model,
			model_parameters: load_model(data_location),
		}

		modelsPool = append(modelsPool, &fetched_model_info)
	}
}

func StartAnomalyDetectionPipeline() {

	utils.WriteToLogs("[PIPELINES]", "Anomaly detection pipeline started")

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	fetchModels()

	for {
		select {
		case <-ticker.C:
			for i := range modelsPool {
				feed_data(modelsPool[i])
			}
		case <-refreshModels:
			utils.WriteToLogs("[PIPELINES]", "Refreshing models pool")
			fetchModels()
		}
	}

}
