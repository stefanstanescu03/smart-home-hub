package pipelines

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type distanceBasedDetector struct {
	mean_distance float32
	std_distance  float32
	k             float32
	lr            float32
	n             uint
	centroid      []float32
	m2_distance   float32
}

type timeseriesWindow struct {
	window      []float32
	numFeatures int
}

type dataset struct {
	data []timeseriesWindow
	size int
}

func csv_to_dataset(filename string, window_size int, param string) dataset {

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Scan()
	line := scanner.Text()
	parts := strings.Split(line, ",")

	var index int

	for i := range parts {
		if parts[i] == param {
			index = i
			break
		}
	}

	var timeserie []float32

	for scanner.Scan() {

		line := scanner.Text()
		parts := strings.Split(line, ",")

		value, err := strconv.ParseFloat(parts[index], 32)
		if err != nil {
			log.Fatal(err)
		}
		final_value := float32(value)

		timeserie = append(timeserie, final_value)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Create data windows
	var data []timeseriesWindow

	for i := 0; i <= len(timeserie)-window_size; i++ {
		window := make([]float32, window_size)
		copy(window, timeserie[i:i+window_size])
		data = append(data, timeseriesWindow{
			window:      window,
			numFeatures: window_size,
		})
	}

	return dataset{
		data: data,
		size: len(data),
	}

}

func mean_centroid(windows []timeseriesWindow) []float32 {

	num_windows := len(windows)
	num_features := windows[0].numFeatures

	var centroid []float32
	centroid = make([]float32, num_features)

	for i := 0; i < num_windows; i++ {
		for j := 0; j < num_features; j++ {
			centroid[j] += windows[i].window[j]
		}
	}

	for i := 0; i < num_features; i++ {
		centroid[i] /= float32(num_windows)
	}

	return centroid

}

func euclidean_distance(a []float32, b []float32) float32 {

	var sum float32 = 0
	for i := range a {
		sum += (a[i] - b[i]) * (a[i] - b[i])
	}

	return float32(math.Sqrt(float64(sum)))

}

func mean(arr []float32) float32 {

	var sum float32 = 0
	for i := range arr {
		sum += arr[i]
	}

	return sum / float32(len(arr))

}

func std(arr []float32) float32 {
	var sum float32
	mean := mean(arr)
	for i := range arr {
		diff := arr[i] - mean
		sum += diff * diff
	}
	return float32(math.Sqrt(float64(sum / float32(len(arr)-1))))
}

func newAnomalyDetectionModel(k float32, lr float32, numFeatures int) *distanceBasedDetector {

	return &distanceBasedDetector{
		k:             k,
		lr:            lr,
		n:             0,
		mean_distance: -1,
		std_distance:  -1,
		centroid:      make([]float32, numFeatures),
		m2_distance:   0,
	}
}

func fit(model *distanceBasedDetector, train_set *dataset) {

	model.centroid = mean_centroid(train_set.data)

	var distances []float32

	for i := range train_set.data {
		distances = append(distances, euclidean_distance(model.centroid, train_set.data[i].window))
	}

	model.mean_distance = mean(distances)
	model.std_distance = std(distances)

	model.n = uint(len(distances))

	for i := range distances {
		model.m2_distance += (distances[i] - model.mean_distance) * (distances[i] - model.mean_distance)
	}
}

func predict(model *distanceBasedDetector, window timeseriesWindow) bool {

	distance := euclidean_distance(model.centroid, window.window)
	threshold := model.mean_distance + model.k*model.std_distance

	if distance > threshold {
		return true
	} else {

		// Adapt center

		for i := range model.centroid {
			model.centroid[i] += model.lr * (window.window[i] - model.centroid[i])
		}

		model.n += 1

		delta := distance - model.mean_distance
		model.mean_distance += delta / float32(model.n)
		model.m2_distance += delta * (distance - model.mean_distance)

		model.std_distance = float32(math.Sqrt(float64(model.m2_distance) / (float64(model.n) - 1)))

		return false
	}
}

// Note that you should never have empty datasets and windows smaller than 2
// Always when a model is initiated collect a lot of data first and after that train
// Model assumes that the training data is stationary and don't have drifts
// This is much fitted for rare anomalies and eventual slow drifts
func StartAnomalyDetectionPipeline() {

	model := newAnomalyDetectionModel(3, 1e-2, 10)

	dataset := csv_to_dataset("/home/stefan/Documents/Projects/licenta/tests/telemetry/test_sensor.csv", 10, "temperature[C]")

	fit(model, &dataset)

	fmt.Println(model)

	testWindow := timeseriesWindow{
		window:      []float32{22, 22, 22, 22, 22, 22, 22, 22, 22, 22},
		numFeatures: 10,
	}

	fmt.Println(predict(model, testWindow))

	fmt.Println(model)

}
