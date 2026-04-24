package pipelines

import (
	"backend/initializers"
	"backend/models"
	"backend/utils"
	"encoding/csv"
	"encoding/gob"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/e-XpertSolutions/go-iforest/v2/iforest"
)

type PersistentModel struct {
	F    *iforest.Forest
	Mean float64
	Std  float64
}

type ModelInfo struct {
	model_metadata *models.AnomalyModel
	model          *PersistentModel
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

func eliminate_units(str string) string {
	if i := strings.Index(str, "["); i != -1 {
		return str[:i]
	} else {
		return str
	}
}

func loadData(filename string, n int, w int, param string) ([][]float64, error) {

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	var index int
	first_line, _ := reader.Read()
	for i := range first_line {
		if eliminate_units(first_line[i]) == param {
			index = i
			break
		}
	}

	var allData []float64
	count := 0

	for count < n {

		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		float_record, _ := strconv.ParseFloat(record[index], 64)
		allData = append(allData, float_record)
		count++
	}

	if len(allData) < w {
		return nil, fmt.Errorf("not enough data to create a window of size %d", w)
	}

	var windows [][]float64
	for i := 0; i <= len(allData)-w; i++ {
		var window []float64
		for j := 0; j < w; j++ {
			window = append(window, allData[i+j])
		}
		windows = append(windows, window)
	}

	return windows, nil

}

func CalculateZScoreStats(data [][]float64) (float64, float64) {
	var sum, count float64
	for _, row := range data {
		for _, val := range row {
			sum += val
			count++
		}
	}
	mean := sum / count

	var sqDiffSum float64
	for _, row := range data {
		for _, val := range row {
			sqDiffSum += math.Pow(val-mean, 2)
		}
	}
	stdDev := math.Sqrt(sqDiffSum / count)

	if stdDev == 0 {
		stdDev = 1
	}

	return mean, stdDev
}

func ApplyZScore(data [][]float64, mean, stdDev float64) {
	for i := range data {
		for j := range data[i] {
			data[i][j] = (data[i][j] - mean) / stdDev
		}
	}
}

func create_model(data_path string, window_size int, param string) *PersistentModel {

	trainData, err := loadData(data_path, 500, window_size, param)
	if err != nil {
		utils.WriteToLogs("ALERTS-HANDLER", fmt.Sprintf("Failed to load data from %s: %v", data_path, err))
	}

	if len(trainData) == 0 {
		utils.WriteToLogs("ALERTS-HANDLER", fmt.Sprintf("Empty train data in %s", data_path))
	}

	var model PersistentModel

	model.Mean, model.Std = CalculateZScoreStats(trainData)

	ApplyZScore(trainData, model.Mean, model.Std)

	f := iforest.NewForest(100, 256, 0.01)
	f.Train(trainData)
	f.Test(trainData)

	model.F = f

	return &model
}

func save_model(filename string, model *PersistentModel) error {

	gob.Register(PersistentModel{})
	gob.Register(iforest.Forest{})

	if model == nil {
		return fmt.Errorf("cannot save a nil model")
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(model)
	if err != nil {
		return fmt.Errorf("gob encode error: %w", err)
	}

	return nil

}

func load_model(filename string) *PersistentModel {

	file, err := os.Open(filename)
	if err != nil {
		utils.WriteToLogs("[PIPELINES]", fmt.Sprintf("File error: %v", err))
		return nil
	}
	defer file.Close()

	var loaded PersistentModel
	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&loaded)
	if err != nil {
		utils.WriteToLogs("[PIPELINES]", fmt.Sprintf("Decode error: %v", err))
		return nil
	}
	return &loaded
}

// FitAndSave creates the initial model and saves it to the desired file.
// This should be called from the controller.
func FitAndSave(filename string, data_path string, param string) {

	model := create_model(data_path, 10, param)
	save_model(filename, model)
	utils.WriteToLogs("[PIPELINES]", fmt.Sprintf("New model created at: %s", filename))

}

func feed_data(model *ModelInfo) {

	if model.model == nil || model.model.F == nil {
		utils.WriteToLogs("[PIPELINES]", "Model or Forest is nil, skipping prediction")
		return
	}

	window_size := 10

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
		if eliminate_units(keys[i]) == param {
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

	var window []float64

	for i := range lines {
		parts := strings.Split(lines[i], ",")
		if param_id >= len(parts) {
			utils.WriteToLogs("[PIPELINES]", fmt.Sprintf("line %d in %s has too few columns", i+1, data_location))
			return
		}

		value, err := strconv.ParseFloat(strings.TrimSpace(parts[param_id]), 64)
		if err != nil {
			utils.WriteToLogs("[PIPELINES]", fmt.Sprintf("Failed to parse value %q at line %d in %s: %v", parts[param_id], i+1, data_location, err))
			return
		}

		window = append(window, (float64(value)-model.model.Mean)/model.model.Std)
	}

	labels, _, err := model.model.F.Predict([][]float64{window})
	if err != nil {
		utils.WriteToLogs("[PIPELINES]", fmt.Sprintf("Prediction error: %v", err))
		return
	}

	if labels[0] == 1 {
		if !alertStates[model.model_metadata.ID] {
			utils.WriteAlert(model.model_metadata.DeviceId, "Anomaly detected", fmt.Sprintf("Unusual behaviour with parameter: %s", param))
			if model.model_metadata.NotifyEmail {
				handleSendEmail(*model.model_metadata)
			}
			alertStates[model.model_metadata.ID] = true
		}
	} else {
		alertStates[model.model_metadata.ID] = false
	}

}

func fetchModels() {

	modelsPool = nil

	alertStates = make(map[uint]bool)

	var models_metadata []models.AnomalyModel
	initializers.DB.Find(&models_metadata)

	for i := range models_metadata {
		fetched_model_info := ModelInfo{
			model_metadata: &models_metadata[i],
			model:          load_model(models_metadata[i].Location),
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
