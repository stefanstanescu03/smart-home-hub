package controllers

import (
	"backend/initializers"
	"backend/models"
	"backend/sockets"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddWidget(c *gin.Context) {

	currUser, _ := c.Get("user")

	var body struct {
		Widgettype  string
		DeviceId    uint
		DashboardId uint
		Payload     string
		Label       string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	var dashboard models.Dashboard
	initializers.DB.First(&dashboard, "user_id = ? and id = ?", currUser.(models.User).ID, body.DashboardId)

	if dashboard.ID == 0 {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You don't have access to this dashboard",
		})
		return
	}

	var device models.Device
	initializers.DB.First(&device, "(user_id = ? or visibility = true) and id = ?", currUser.(models.User).ID, body.DeviceId)

	if device.ID == 0 {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You don't have access to this device",
		})
		return
	}

	widget := models.Widget{
		Widgettype:  body.Widgettype,
		DeviceId:    body.DeviceId,
		DashboardId: body.DashboardId,
		Payload:     body.Payload,
		Label:       body.Label,
	}

	res := initializers.DB.Create(&widget)

	if res.Error != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"widget": widget,
	})

}

func GetWidgets(c *gin.Context) {

	currUser, _ := c.Get("user")
	dashboardId := c.Param("id")

	var dashboard models.Dashboard
	initializers.DB.First(&dashboard, "user_id = ? and id = ?", currUser.(models.User).ID, dashboardId)

	if dashboard.ID == 0 {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You don't have access to this dashboard",
		})
		return
	}

	var widgets []models.Widget
	initializers.DB.Find(&widgets, "dashboard_id = ?", dashboardId)

	c.JSON(http.StatusOK, gin.H{
		"widgets": widgets,
	})
}

func GetPublicWidgets(c *gin.Context) {

	dashboardId := c.Param("id")

	var dashboard models.Dashboard
	initializers.DB.First(&dashboard, "visibility = true and id = ?", dashboardId)

	if dashboard.ID == 0 {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You don't have access to this dashboard",
		})
		return
	}

	var widgets []models.Widget
	initializers.DB.Find(&widgets, "dashboard_id = ?", dashboardId)

	c.JSON(http.StatusOK, gin.H{
		"widgets": widgets,
	})
}

func DeleteWidget(c *gin.Context) {

	currUser, _ := c.Get("user")
	id := c.Param("id")

	var widget models.Widget
	initializers.DB.First(&widget, "id = ?", id)
	var dashboard models.Dashboard
	initializers.DB.First(&dashboard, "user_id = ? and id = ?", currUser.(models.User).ID, widget.DashboardId)

	if dashboard.ID == 0 {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You don't have access to this dashboard",
		})
		return
	}

	initializers.DB.Unscoped().Delete(models.Widget{}, "id = ?", id)

	c.JSON(http.StatusOK, gin.H{
		"message": "widget deleted",
	})

}

func SendCommand(c *gin.Context) {

	currUser, _ := c.Get("user")
	var body struct {
		Widget   uint
		DeviceId uint
		Payload  string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	var widget models.Widget
	initializers.DB.First(&widget, "device_id = ? and id = ?", body.DeviceId, body.Widget)
	if widget.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Widget is not linked with this device",
		})
		return
	}

	var device models.Device
	initializers.DB.First(&device, "user_id = ? and id = ?", currUser.(models.User).ID, body.DeviceId)
	if device.ID == 0 {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You don't have access to this device",
		})
		return
	}

	sockets.PublishCmd(device.Ident, body.Payload)

	c.JSON(http.StatusOK, gin.H{
		"ident":   device.Ident,
		"payload": body.Payload,
	})
}
