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
	N             uint
	Centroid      []float32
	M2_distance   float32
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
	var data []TimeseriesWindow

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
	return float32(math.Sqrt(float64(sum / float32(len(arr)-1))))
}

func NewAnomalyDetectionModel(k float32, lr float32, numFeatures int) *DistanceBasedDetector {

	return &DistanceBasedDetector{
		K:             k,
		Lr:            lr,
		N:             0,
		Mean_distance: -1,
		Std_distance:  -1,
		Centroid:      make([]float32, numFeatures),
		M2_distance:   0,
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

	model.N = uint(len(distances))

	model.M2_distance = model.Std_distance * model.Std_distance * float32(model.N-1)
}

func Predict(model *DistanceBasedDetector, window TimeseriesWindow) bool {

	distance := euclidean_distance(model.Centroid, window.Window)
	threshold := model.Mean_distance + model.K*model.Std_distance

	if distance > threshold {
		return true
	} else {

		// Adapt center

		for i := range model.Centroid {
			model.Centroid[i] += model.Lr * (window.Window[i] - model.Centroid[i])
		}

		adaptedDistance := euclidean_distance(model.Centroid, window.Window)

		model.N += 1

		delta := adaptedDistance - model.Mean_distance
		model.Mean_distance += delta / float32(model.N)
		model.M2_distance += delta * (adaptedDistance - model.Mean_distance)

		model.Std_distance = float32(math.Sqrt(float64(model.M2_distance) / (float64(model.N) - 1)))

		return false
	}
}
