package services

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAlertRoutes(r *gin.Engine) {
	alertRoutes := r.Group("/alert")
	{
		alertRoutes.POST("/create", middleware.RequireAuth, controllers.AddAlert)
		alertRoutes.PUT("/update/:id", middleware.RequireAuth, controllers.UpdateAlert)
		alertRoutes.GET("/:id", middleware.RequireAuth, controllers.GetAlert)
		alertRoutes.GET("/", middleware.RequireAuth, controllers.GetAlerts)
		alertRoutes.DELETE("/delete/:id", middleware.RequireAuth, controllers.DeleteAlert)
	}
}
