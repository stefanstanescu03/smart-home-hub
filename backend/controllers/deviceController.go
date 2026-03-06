package controllers

import (
	"backend/initializers"
	"backend/models"
	"backend/pipelines"
	"backend/sockets"
	"backend/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func AddDevice(c *gin.Context) {

	currUser, _ := c.Get("user")

	var body struct {
		Name  string
		Ident string
		// Csv_location string
		Visibility bool
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	device := models.Device{
		Name:         body.Name,
		Ident:        body.Ident,
		Csv_location: os.Getenv("DATA_LOCATION") + "/" + utils.Format(body.Name) + fmt.Sprint(currUser.(models.User).ID) + ".csv",
		Visibility:   body.Visibility,
		UserId:       currUser.(models.User).ID,
	}

	res := initializers.DB.Create(&device)

	if res.Error != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"device": device,
	})

}

func UpdateDevice(c *gin.Context) {

	currUser, _ := c.Get("user")
	id := c.Param("id")

	var body struct {
		Name         string
		Ident        string
		Csv_location string
		Visibility   bool
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	var device models.Device
	initializers.DB.First(&device, "id = ? and user_id = ?", id, currUser.(models.User).ID)

	if device.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Device not found",
		})
		return
	}

	device.Name = body.Name
	device.Ident = body.Ident
	device.Csv_location = body.Csv_location
	device.Visibility = body.Visibility

	initializers.DB.Save(&device)

	c.JSON(http.StatusOK, gin.H{
		"device": device,
	})
}

func GetDevice(c *gin.Context) {

	currUser, _ := c.Get("user")
	id := c.Param("id")

	var device models.Device
	initializers.DB.First(&device, "id = ? and (user_id = ? or visibility = true)", id, currUser.(models.User).ID)

	if device.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Device not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"device": device,
	})

}

func GetDevices(c *gin.Context) {

	currUser, _ := c.Get("user")
	var devices []models.Device

	initializers.DB.Find(&devices, "user_id = ? or visibility = true", currUser.(models.User).ID)

	c.JSON(http.StatusOK, gin.H{
		"devices": devices,
	})
}

func GetAllDevices(c *gin.Context) {
	currUser, _ := c.Get("user")

	// For admins only
	if currUser.(models.User).Role != "admin" {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	type returnedDevice struct {
		ID         uint
		Name       string
		Ident      string
		Visibility bool
		Owner      string
	}

	var devices []models.Device
	initializers.DB.Find(&devices)

	var returnedDevices []returnedDevice

	for i := range devices {

		var user models.User
		initializers.DB.First(&user, "ID = ?", devices[i].UserId)

		device := returnedDevice{
			ID:         devices[i].ID,
			Name:       devices[i].Name,
			Ident:      devices[i].Ident,
			Visibility: devices[i].Visibility,
			Owner:      user.Username,
		}

		returnedDevices = append(returnedDevices, device)
	}

	c.JSON(http.StatusOK, gin.H{
		"devices": returnedDevices,
	})
}

func GetUserDevices(c *gin.Context) {
	currUser, _ := c.Get("user")
	var devices []models.Device

	initializers.DB.Find(&devices, "user_id = ?", currUser.(models.User).ID)

	c.JSON(http.StatusOK, gin.H{
		"devices": devices,
	})
}

func DeleteDevice(c *gin.Context) {

	currUser, _ := c.Get("user")
	id := c.Param("id")

	if currUser.(models.User).Role != "admin" {
		var device models.Device
		initializers.DB.First(&device, "id = ? and user_id = ?", id, currUser.(models.User).ID)
		if device.ID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Device not found",
			})
			return
		}
	}

	initializers.DB.Unscoped().Delete(&models.Alert{}, "device_id = ?", id)
	initializers.DB.Unscoped().Delete(&models.Widget{}, "device_id = ?", id)
	initializers.DB.Unscoped().Delete(&models.AnomalyModel{}, "device_id = ?", id)
	initializers.DB.Unscoped().Delete(&models.Device{}, "id = ?", id)

	sockets.NotifyAlertsHandler()
	pipelines.NotifyAnomalyPipeline()

	c.JSON(http.StatusOK, gin.H{
		"message": "device deleted",
	})
}

func IsDeviceConnected(c *gin.Context) {

	ident := c.Param("ident")
	_, ok := sockets.ConnectionPool.Load(ident)

	if ok {
		c.JSON(http.StatusOK, gin.H{
			"status": true,
		})
		return
	}

	// Maybe we try to check by id

	var device models.Device
	initializers.DB.First(&device, "id = ?", ident)
	if device.ID == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
		})
		return
	}

	_, ok = sockets.ConnectionPool.Load(device.Ident)

	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": true,
		})
	}
}

func GetPublicDevices(c *gin.Context) {

	var devices []models.Device

	initializers.DB.Find(&devices, "visibility = ?", 1)

	c.JSON(http.StatusOK, gin.H{
		"devices": devices,
	})
}

func GetPublicDevice(c *gin.Context) {

	id := c.Param("id")

	var device models.Device
	initializers.DB.Find(&device, "visibility = true and id = ?", id)

	if device.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Device not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"device": device,
	})
}

func GetDeviceState(c *gin.Context) {

	id := c.Param("id")

	var device models.Device
	initializers.DB.Find(&device, "visibility = true and id = ?", id)

	state, ok := sockets.DeviceStates.Load(device.Ident)

	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"state": "not registered",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"state": state,
		})
	}

}

func GetDeviceParams(c *gin.Context) {

	currUser, _ := c.Get("user")
	id := c.Param("id")

	var device models.Device
	initializers.DB.Find(&device, "user_id = ? and id = ?", currUser.(models.User).ID, id)

	if device.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Device not found",
		})
		return
	}

	first_line := utils.GetFirstLine(device.Csv_location)
	c.JSON(http.StatusOK, gin.H{
		"params": first_line,
	})

}

func GetDeviceData(c *gin.Context) {

	currUser, _ := c.Get("user")
	id := c.Param("id")

	path := c.Query("path")
	num := c.Query("num")
	param := c.Query("param")

	var device models.Device
	initializers.DB.Find(&device, "user_id = ? and id = ?", currUser.(models.User).ID, id)

	if device.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Device not found",
		})
		return
	}

	splits := strings.Split(device.Csv_location, "/")
	filename := splits[len(splits)-1]

	api := fmt.Sprintf("http://localhost:8000/%s/?filename=%s&num=%s&param=%s", path, filename, num, param)

	type apiResponse struct {
		Values []map[string]interface{} `json:"values"`
	}

	resp, err := http.Get(api)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to Fastapi service"})
		return
	}

	var data apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse Fastapi response"})
		return
	}

	c.JSON(http.StatusOK, data)

}

func GetDeviceForecast(c *gin.Context) {

	currUser, _ := c.Get("user")
	id := c.Param("id")

	path := c.Query("path")
	lag := c.Query("lag")
	steps := c.Query("steps")
	param := c.Query("param")

	var device models.Device
	initializers.DB.Find(&device, "user_id = ? and id = ?", currUser.(models.User).ID, id)

	if device.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Device not found",
		})
		return
	}

	splits := strings.Split(device.Csv_location, "/")
	filename := splits[len(splits)-1]

	api := fmt.Sprintf("http://localhost:8000/predict/%s/?filename=%s&lag=%s&steps=%s&param=%s", path, filename, lag, steps, param)

	type apiResponse struct {
		Values []float32 `json:"values"`
	}

	resp, err := http.Get(api)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to Fastapi service"})
		return
	}

	var data apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse Fastapi response"})
		return
	}

	c.JSON(http.StatusOK, data)

}
