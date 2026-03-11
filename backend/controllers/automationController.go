package controllers

import (
	"backend/initializers"
	"backend/models"
	"backend/sockets"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddAutomation(c *gin.Context) {

	currUser, _ := c.Get("user")

	var body struct {
		Name          string
		Device1Id     uint
		Device2Id     uint
		Payload       string
		Condition     string
		ScheduleStart string
		ScheduleEnd   string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	automation := models.Automation{
		Name:          body.Name,
		Device1Id:     body.Device1Id,
		Device2Id:     body.Device2Id,
		Payload:       body.Payload,
		Condition:     body.Condition,
		UserId:        currUser.(models.User).ID,
		ScheduleStart: body.ScheduleStart,
		ScheduleEnd:   body.ScheduleEnd,
	}

	res := initializers.DB.Create(&automation)

	if res.Error != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	sockets.NotifyAutomationsHandler()

	c.JSON(http.StatusOK, gin.H{
		"automation": automation,
	})

}

func GetAutomation(c *gin.Context) {

	currUser, _ := c.Get("user")
	id := c.Param("id")

	var automation models.Automation
	initializers.DB.First(&automation, "id = ? and user_id = ?", id, currUser.(models.User).ID)

	if automation.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Automation not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"automation": automation,
	})

}

func GetAutomations(c *gin.Context) {

	currUser, _ := c.Get("user")

	var automations []models.Automation
	initializers.DB.Find(&automations, "user_id = ?", currUser.(models.User).ID)

	c.JSON(http.StatusOK, gin.H{
		"automations": automations,
	})
}

func DeleteAutomation(c *gin.Context) {

	currUser, _ := c.Get("user")
	id := c.Param("id")

	initializers.DB.Unscoped().Delete(&models.Automation{}, "id = ? and user_id = ?", id, currUser.(models.User).ID)

	sockets.NotifyAutomationsHandler()

	c.JSON(http.StatusOK, gin.H{
		"message": "automation deleted",
	})
}
