package services

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAutomationRoutes(r *gin.Engine) {
	automationRoutes := r.Group("/api/automation")
	{
		automationRoutes.POST("/create", middleware.RequireAuth, controllers.AddAutomation)
		automationRoutes.GET("/:id", middleware.RequireAuth, controllers.GetAutomation)
		automationRoutes.GET("/", middleware.RequireAuth, controllers.GetAutomations)
		automationRoutes.DELETE("/delete/:id", middleware.RequireAuth, controllers.DeleteAutomation)
	}
}
