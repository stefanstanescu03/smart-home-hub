package services

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAlertRoutes(r *gin.Engine) {
	alertRoutes := r.Group("/api/alert")
	{
		alertRoutes.POST("/create", middleware.RequireAuth, controllers.AddAlert)
		alertRoutes.PUT("/update/:id", middleware.RequireAuth, controllers.UpdateAlert)
		alertRoutes.GET("/:device_id/:id", middleware.RequireAuth, controllers.GetAlert)
		alertRoutes.GET("/:device_id", middleware.RequireAuth, controllers.GetAlerts)
		alertRoutes.DELETE("/delete/:id", middleware.RequireAuth, controllers.DeleteAlert)
	}
}
