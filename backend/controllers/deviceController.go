package controllers

import (
	"backend/initializers"
	"backend/models"
	"backend/sockets"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddDevice(c *gin.Context) {

	currUser, _ := c.Get("user")

	var body struct {
		Name         string
		Ip           string
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
		Ip:           body.Ip,
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
		Ip           string
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
	device.Ip = body.Ip
	device.Csv_location = body.Csv_location
	device.Visibility = body.Visibility

	initializers.DB.Save(&device)

	c.JSON(http.StatusOK, gin.H{
		"device": device,
	})
}

func GetDevice(c *gin.Context) {

	currUser, _ := c.Get("user")
	name := c.Param("name")

	var device models.Device
	initializers.DB.First(&device, "name = ? and (user_id = ? or visibility = true)", name, currUser.(models.User).ID)

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

	initializers.DB.Find(&devices, "user_id = ?", currUser.(models.User).ID)

	c.JSON(http.StatusOK, gin.H{
		"devices": devices,
	})
}

func DeleteDevice(c *gin.Context) {

	currUser, _ := c.Get("user")
	id := c.Param("id")

	initializers.DB.Unscoped().Delete(&models.Device{}, "id = ? and user_id = ?", id, currUser.(models.User).ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "device deleted",
	})
}

func IsDeviceConnected(c *gin.Context) {

	ip := c.Param("ip")
	value, ok := sockets.TelemetryConnectionPool.Load(ip)

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

func GetDiscoveredDevices(c *gin.Context) {

	devices := sockets.AcknowledgedDevices

	c.JSON(http.StatusOK, gin.H{
		"devices": devices,
	})
}

func GetPublicDevices(c *gin.Context) {

	var devices []models.Device

	initializers.DB.Find(&devices, "visibility = ?", 1)

	c.JSON(http.StatusOK, gin.H{
		"devices": devices,
	})
}
