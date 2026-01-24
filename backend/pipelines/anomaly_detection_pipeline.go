package pipelines

import (
	"backend/pipelines/core"
	"encoding/json"
	"log"
	"os"
)

var ModelsPool []*core.DistanceBasedDetector

func create_model(data_path string, window_size int, param string) {

	model := core.NewAnomalyDetectionModel(3, 1e-2, window_size)
	dataset := core.Csv_to_dataset(data_path, window_size, param)
	core.Fit(model, &dataset)

	ModelsPool = append(ModelsPool, model)
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

func load_model(filename string) {

	json_value, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	model := core.NewAnomalyDetectionModel(3, 1e-2, 10)
	json.Unmarshal(json_value, model)

	ModelsPool = append(ModelsPool, model)

}

func StartAnomalyDetectionPipeline() {

	// create_model("/home/stefan/Documents/Projects/licenta/tests/telemetry/test_sensor.csv", 10, "temperature[C]")
	// create_model("/home/stefan/Documents/Projects/licenta/tests/telemetry/test_sensor.csv", 10, "humidity[%]")

	// save_model("/home/stefan/Documents/Projects/licenta/backend/ml-models/temperature_model.json", Models[0])

	load_model("/home/stefan/Documents/Projects/licenta/backend/ml-models/temperature_model.json")

}
