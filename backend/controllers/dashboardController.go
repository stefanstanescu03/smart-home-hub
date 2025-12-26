package controllers

import (
	"backend/initializers"
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddDashboard(c *gin.Context) {

	currUser, _ := c.Get("user")

	var body struct {
		Name       string
		Visibility bool
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	dashboard := models.Dashboard{
		Name:       body.Name,
		Visibility: body.Visibility,
		UserId:     currUser.(models.User).ID,
	}

	res := initializers.DB.Create(&dashboard)

	if res.Error != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"dashboard": dashboard,
	})
}

func UpdateDashboard(c *gin.Context) {

	currUser, _ := c.Get("user")
	id := c.Param("id")

	var body struct {
		Name       string
		Visibility bool
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	var dashboard models.Dashboard
	initializers.DB.First(&dashboard, "user_id = ? and id = ?", currUser.(models.User).ID, id)

	if dashboard.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Dashboard not found",
		})
		return
	}

	dashboard.Name = body.Name
	dashboard.Visibility = body.Visibility

	initializers.DB.Save(&dashboard)

	c.JSON(http.StatusOK, gin.H{
		"dashboard": dashboard,
	})

}

func DeleteDashboard(c *gin.Context) {

	currUser, _ := c.Get("user")
	id := c.Param("id")

	var dashboard models.Dashboard
	initializers.DB.First(&dashboard, "user_id = ? and id = ?", currUser.(models.User).ID, id)

	if dashboard.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Dashboard not found",
		})
		return
	}

	initializers.DB.Unscoped().Delete(models.Widget{}, "dashboard_id = ?", id)
	initializers.DB.Unscoped().Delete(models.Dashboard{}, "id = ? and user_id = ?", id, currUser.(models.User).ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "dashboard deleted",
	})

}

func GetDashboard(c *gin.Context) {

	currUser, _ := c.Get("user")
	id := c.Param("id")

	var dashboard models.Dashboard
	initializers.DB.First(&dashboard, "(user_id = ? or visibility = true) and id = ?", currUser.(models.User).ID, id)

	if dashboard.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Dashboard not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"dashboard": dashboard,
	})

}

func GetDashboards(c *gin.Context) {

	currUser, _ := c.Get("user")

	var dashboards []models.Dashboard
	initializers.DB.Find(&dashboards, "user_id = ?", currUser.(models.User).ID)

	c.JSON(http.StatusOK, gin.H{
		"dashboards": dashboards,
	})
}

func GetPublicDashboards(c *gin.Context) {
	var dashboards []models.Dashboard
	initializers.DB.Find(&dashboards, "visibility = true")

	c.JSON(http.StatusOK, gin.H{
		"dashboards": dashboards,
	})
}

func GetDashboardPublic(c *gin.Context) {

	id := c.Param("id")

	var dashboard models.Dashboard
	initializers.DB.First(&dashboard, "visibility = true and id = ?", id)

	if dashboard.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Dashboard not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"dashboard": dashboard,
	})

}
