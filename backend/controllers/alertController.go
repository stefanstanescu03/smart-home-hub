package controllers

import (
	"backend/initializers"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddAlert(c *gin.Context) {

	currUser, _ := c.Get("user")

	var body struct {
		Subject   string
		Content   string
		Condition string
		DeviceId  uint
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

	alert := models.Alert{
		Subject:   body.Subject,
		Content:   body.Content,
		Condition: body.Condition,
		UserId:    currUser.(models.User).ID,
		DeviceId:  body.DeviceId,
	}

	res := initializers.DB.Create(&alert)

	if res.Error != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"alert": alert,
	})
}

func UpdateAlert(c *gin.Context) {

	currUser, _ := c.Get("user")
	id := c.Param("id")

	var body struct {
		Subject   string
		Content   string
		Condition string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	var alert models.Alert
	initializers.DB.First(&alert, "id = ? and user_id = ?", id, currUser.(models.User).ID)

	if alert.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "alert not found",
		})
		return
	}

	alert.Subject = body.Subject
	alert.Content = body.Content
	alert.Condition = body.Condition

	initializers.DB.Save(&alert)

	c.JSON(http.StatusOK, gin.H{
		"alert": alert,
	})
}

func GetAlert(c *gin.Context) {

	currUser, _ := c.Get("user")
	id := c.Param("id")

	var alert models.Alert
	initializers.DB.First(&alert, "id = ? and user_id = ?", id, currUser.(models.User).ID)

	if alert.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "alert not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"alert": alert,
	})

}

func GetAlerts(c *gin.Context) {

	currUser, _ := c.Get("user")
	device_id := c.Param("device_id")

	var alerts []models.Alert
	initializers.DB.Find(&alerts, "user_id = ? and device_id = ?", currUser.(models.User).ID, device_id)

	c.JSON(http.StatusOK, gin.H{
		"alerts": alerts,
	})
}

func DeleteAlert(c *gin.Context) {

	currUser, _ := c.Get("user")
	id := c.Param("id")

	initializers.DB.Unscoped().Delete(&models.Alert{}, "id = ? and user_id = ?", id, currUser.(models.User).ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "alert deleted",
	})
}
