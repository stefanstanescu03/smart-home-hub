package core

import (
	"bufio"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type DistanceBasedDetector struct {
	Mean_distance float32
	Std_distance  float32
	K             float32
	Lr            float32
	Centroid      []float32
}

type TimeseriesWindow struct {
	Window      []float32
	NumFeatures int
}

type Dataset struct {
	Data []TimeseriesWindow
	Size int
}

func Csv_to_dataset(filename string, window_size int, param string) Dataset {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("failed to open file %q: %v", filename, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Parse header
	if !scanner.Scan() {
		log.Fatal("CSV file is empty or missing header row")
	}
	headers := strings.Split(scanner.Text(), ",")

	index := -1
	for i, h := range headers {
		if strings.TrimSpace(h) == param {
			index = i
			break
		}
	}
	if index == -1 {
		log.Fatalf("column %q not found in CSV header", param)
	}

	// Parse rows
	var timeserie []float32
	lineNum := 1
	for scanner.Scan() {
		lineNum++
		parts := strings.Split(scanner.Text(), ",")
		if index >= len(parts) {
			log.Fatalf("line %d: expected at least %d columns, got %d", lineNum, index+1, len(parts))
		}
		value, err := strconv.ParseFloat(strings.TrimSpace(parts[index]), 32)
		if err != nil {
			log.Fatalf("line %d: could not parse value %q as float: %v", lineNum, parts[index], err)
		}
		timeserie = append(timeserie, float32(value))
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("error reading file: %v", err)
	}

	if len(timeserie) < window_size {
		log.Fatalf("not enough data: need at least %d rows, got %d", window_size, len(timeserie))
	}

	// Build sliding windows
	data := make([]TimeseriesWindow, 0, len(timeserie)-window_size+1)
	for i := 0; i <= len(timeserie)-window_size; i++ {
		window := make([]float32, window_size)
		copy(window, timeserie[i:i+window_size])
		data = append(data, TimeseriesWindow{
			Window:      window,
			NumFeatures: window_size,
		})
	}

	return Dataset{
		Data: data,
		Size: len(data),
	}
}

func mean_centroid(windows []TimeseriesWindow) []float32 {

	num_windows := len(windows)
	num_features := windows[0].NumFeatures

	var centroid []float32
	centroid = make([]float32, num_features)

	for i := 0; i < num_windows; i++ {
		for j := 0; j < num_features; j++ {
			centroid[j] += windows[i].Window[j]
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
	return float32(math.Sqrt(float64(sum / float32(len(arr)))))
}

func NewAnomalyDetectionModel(k float32, lr float32, numFeatures int) *DistanceBasedDetector {

	return &DistanceBasedDetector{
		K:             k,
		Lr:            lr,
		Mean_distance: -1,
		Std_distance:  -1,
		Centroid:      make([]float32, numFeatures),
	}
}

func Fit(model *DistanceBasedDetector, train_set *Dataset) {

	model.Centroid = mean_centroid(train_set.Data)

	var distances []float32

	for i := range train_set.Data {
		distances = append(distances, euclidean_distance(model.Centroid, train_set.Data[i].Window))
	}

	model.Mean_distance = mean(distances)
	model.Std_distance = std(distances)
}

func Predict(model *DistanceBasedDetector, window TimeseriesWindow) bool {

	if model.Mean_distance < 0 || model.Std_distance < 0 {
		log.Fatal("model must be fitted before calling Predict")
	}

	distance := euclidean_distance(model.Centroid, window.Window)
	threshold := model.Mean_distance + model.K*model.Std_distance

	if distance > threshold {
		return true
	} else {

		// Adapt center

		for i := range model.Centroid {
			model.Centroid[i] += model.Lr * (window.Window[i] - model.Centroid[i])
		}
		distance = euclidean_distance(model.Centroid, window.Window)
		diff := distance - model.Mean_distance
		variance := model.Std_distance * model.Std_distance
		model.Mean_distance += model.Lr * diff
		model.Std_distance = float32(math.Sqrt(float64(variance + model.Lr*(diff*diff-variance))))

		return false
	}
}
