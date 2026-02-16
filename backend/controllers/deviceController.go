package controllers

import (
	"backend/initializers"
	"backend/models"
	"backend/pipelines"
	"backend/sockets"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddDevice(c *gin.Context) {

	currUser, _ := c.Get("user")

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

	device := models.Device{
		Name:         body.Name,
		Ident:        body.Ident,
		Csv_location: body.Csv_location,
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

	var device models.Device
	initializers.DB.First(&device, "id = ? and user_id = ?", id, currUser.(models.User).ID)
	if device.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Device not found",
		})
		return
	}

	initializers.DB.Unscoped().Delete(&models.Alert{}, "device_id = ?", id)
	initializers.DB.Unscoped().Delete(&models.Widget{}, "device_id = ?", id)
	initializers.DB.Unscoped().Delete(&models.AnomalyModel{}, "device_id = ?", id)
	initializers.DB.Unscoped().Delete(&models.Device{}, "id = ? and user_id = ?", id, currUser.(models.User).ID)

	sockets.NotifyAlertsHandler()
	pipelines.NotifyAnomalyPipeline()

	c.JSON(http.StatusOK, gin.H{
		"message": "device deleted",
	})
}

func IsDeviceConnected(c *gin.Context) {

	ident := c.Param("ident")
	value, ok := sockets.TelemetryConnectionPool.Load(ident)

	if !ok {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": value,
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
