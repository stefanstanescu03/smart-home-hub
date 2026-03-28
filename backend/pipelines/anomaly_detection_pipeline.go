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

	model := core.NewAnomalyDetectionModel(3, 1e-2, window_size)
	dataset := core.Csv_to_dataset(data_path, window_size, param)
	core.Fit(model, &dataset)

	return model
}

func save_model(filename string, model *core.DistanceBasedDetector) {

	json_value, err := json.Marshal(model)
	if err != nil {
		log.Fatalf("failed to serialize model to %q: %v", filename, err)
	}

	err = os.WriteFile(filename, json_value, 0644)
	if err != nil {
		log.Fatalf("failed to write model file %q: %v", filename, err)
	}

}

func load_model(filename string) *core.DistanceBasedDetector {

	json_value, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("failed to read model file %q: %v", filename, err)
	}

	var model core.DistanceBasedDetector
	if err := json.Unmarshal(json_value, &model); err != nil {
		log.Fatalf("failed to deserialize model from %q: %v", filename, err)
	}

	return &model
}

// FitAndSave creates the initial model and saves it to the desired file.
// This should be called from the controller.
func FitAndSave(filename string, data_path string, param string) {

	model := create_model(data_path, 10, param)
	save_model(filename, model)
	utils.WriteToLogs("[PIPELINES]", fmt.Sprintf("New model created at: %s", filename))

}

func feed_data(model *ModelInfo) {

	window_size := len(model.model_parameters.Centroid)

	device_id := model.model_metadata.DeviceId
	var device models.Device
	initializers.DB.First(&device, "id = ?", device_id)

	if device.ID == 0 {
		utils.WriteToLogs("[PIPELINES]", fmt.Sprintf("Somehow the device was not found: %d", device_id))
		return
	}

	data_location := device.Csv_location
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
	if param_id == -1 {
		utils.WriteToLogs("[PIPELINES]", fmt.Sprintf("param %q not found in CSV header at %s", param, data_location))
		return
	}

	lines, err := utils.GetLastLines(data_location, window_size)
	if err != nil {
		utils.WriteToLogs("[PIPELINES]", fmt.Sprintf("Failed to get last values in: %s", data_location))
		return
	}

	var window core.TimeseriesWindow
	window.NumFeatures = window_size

	for i := range lines {
		parts := strings.Split(lines[i], ",")
		if param_id >= len(parts) {
			utils.WriteToLogs("[PIPELINES]", fmt.Sprintf("line %d in %s has too few columns", i+1, data_location))
			return
		}

		value, err := strconv.ParseFloat(strings.TrimSpace(parts[param_id]), 32)
		if err != nil {
			utils.WriteToLogs("[PIPELINES]", fmt.Sprintf("Failed to parse value %q at line %d in %s: %v", parts[param_id], i+1, data_location, err))
			return
		}

		window.Window = append(window.Window, float32(value))
	}

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

	alertStates = make(map[uint]bool)

	var models_metadata []models.AnomalyModel
	initializers.DB.Find(&models_metadata)

	for i := range models_metadata {
		fetched_model_info := ModelInfo{
			model_metadata:   &models_metadata[i],
			model_parameters: load_model(models_metadata[i].Location),
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
