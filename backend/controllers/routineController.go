package controllers

import (
	"backend/initializers"
	"backend/models"
	"backend/sockets"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddRoutine(c *gin.Context) {

	currUser, _ := c.Get("user")

	var body struct {
		Name      string
		DeviceId  uint
		Payload   string
		StartTime string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	routine := models.Routine{
		Name:      body.Name,
		DeviceId:  body.DeviceId,
		Payload:   body.Payload,
		UserId:    currUser.(models.User).ID,
		StartTime: body.StartTime,
	}

	res := initializers.DB.Create(&routine)

	if res.Error != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	sockets.NotifyAutomationsHandler()

	c.JSON(http.StatusOK, gin.H{
		"routine": routine,
	})

}

func GetRoutine(c *gin.Context) {

	currUser, _ := c.Get("user")
	id := c.Param("id")

	var routine models.Routine
	initializers.DB.First(&routine, "id = ? and user_id = ?", id, currUser.(models.User).ID)

	if routine.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Routine not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"routine": routine,
	})

}

func GetRoutines(c *gin.Context) {

	currUser, _ := c.Get("user")

	var routines []models.Routine
	initializers.DB.Find(&routines, "user_id = ?", currUser.(models.User).ID)

	c.JSON(http.StatusOK, gin.H{
		"routines": routines,
	})
}

func DeleteRoutine(c *gin.Context) {

	currUser, _ := c.Get("user")
	id := c.Param("id")

	initializers.DB.Unscoped().Delete(&models.Routine{}, "id = ? and user_id = ?", id, currUser.(models.User).ID)

	sockets.NotifyAutomationsHandler()

	c.JSON(http.StatusOK, gin.H{
		"message": "routine deleted",
	})
}
