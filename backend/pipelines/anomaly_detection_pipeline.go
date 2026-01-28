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
	"time"
)

type ModelInfo struct {
	model_metadata   *models.AnomalyModel
	model_parameters *core.DistanceBasedDetector
}

var modelsPool []*ModelInfo

var refreshModels = make(chan struct{}, 1)

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
	model := core.NewAnomalyDetectionModel(3, 1e-2, 10)
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

	for _, model := range modelsPool {
		fmt.Println(model.model_metadata.Location)
	}

	for {
		select {
		case <-ticker.C:
		case <-refreshModels:
			utils.WriteToLogs("[PIPELINES]", "Refreshing models pool")
			fetchModels()
		}
	}

}
