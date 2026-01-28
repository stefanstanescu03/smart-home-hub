package controllers

import (
	"backend/initializers"
	"backend/models"
	"backend/pipelines"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddAnomalyModel(c *gin.Context) {

	currUser, _ := c.Get("user")

	var body struct {
		Location string
		Param    string
		DeviceId uint
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	var device models.Device
	initializers.DB.Find(&device, "user_id = ? and id = ?", currUser.(models.User).ID, body.DeviceId)
	if device.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "device not found",
		})
		return
	}

	model := models.AnomalyModel{
		Location:    body.Location,
		Param:       body.Param,
		NotifyEmail: false,
		UserId:      currUser.(models.User).ID,
		DeviceId:    body.DeviceId,
	}

	res := initializers.DB.Create(&model)

	if res.Error != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	pipelines.FitAndSave(body.Location, device.Csv_location, body.Param)

	pipelines.NotifyAnomalyPipeline()

	c.JSON(http.StatusOK, gin.H{
		"model": model,
	})
}

func UpdateAnomalyModel(c *gin.Context) {

	currUser, _ := c.Get("user")
	id := c.Param("id")

	var body struct {
		Location    string
		NotifyEmail bool
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	var model models.AnomalyModel
	initializers.DB.First(&model, "id = ? and user_id = ?", id, currUser.(models.User).ID)

	if model.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "model not found",
		})
		return
	}

	model.Location = body.Location
	model.NotifyEmail = body.NotifyEmail

	initializers.DB.Save(&model)

	c.JSON(http.StatusOK, gin.H{
		"model": model,
	})
}

func GetAnomalyModel(c *gin.Context) {

	currUser, _ := c.Get("user")
	id := c.Param("id")

	var model models.AnomalyModel
	initializers.DB.First(&model, "id = ? and user_id = ?", id, currUser.(models.User).ID)

	if model.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "model not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"model": model,
	})

}

func GetAnomalyModels(c *gin.Context) {

	currUser, _ := c.Get("user")
	device_id := c.Param("device_id")

	var anomalyModels []models.AnomalyModel
	initializers.DB.Find(&anomalyModels, "user_id = ? and device_id = ?", currUser.(models.User).ID, device_id)

	c.JSON(http.StatusOK, gin.H{
		"models": anomalyModels,
	})
}

func DeleteAnomalyModel(c *gin.Context) {

	currUser, _ := c.Get("user")
	id := c.Param("id")

	initializers.DB.Unscoped().Delete(&models.AnomalyModel{}, "id = ? and user_id = ?", id, currUser.(models.User).ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "model deleted",
	})
}
