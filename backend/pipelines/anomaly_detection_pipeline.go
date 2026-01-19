package pipelines

import (
	"backend/pipelines/core"
	"fmt"
)

type modelMetadata struct {
	param         string
	data_path     string
	mean_distance float32
	std_distance  float32
	k             float32
	lr            float32
	n             uint
	centroid      []float32
	m2_distance   float32
}

var Models []*core.DistanceBasedDetector

func create_model(data_path string, window_size int, param string) {

	model := core.NewAnomalyDetectionModel(3, 1e-2, window_size)
	dataset := core.Csv_to_dataset(data_path, window_size, param)
	core.Fit(model, &dataset)

	Models = append(Models, model)
}

func StartAnomalyDetectionPipeline() {

	create_model("/home/stefan/Documents/Projects/licenta/tests/telemetry/test_sensor.csv", 10, "temperature[C]")
	create_model("/home/stefan/Documents/Projects/licenta/tests/telemetry/test_sensor.csv", 10, "humidity[%]")

	for i := range Models {
		fmt.Println(Models[i])
	}

}
