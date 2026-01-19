package services

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterAnomalyModelRoutes(r *gin.Engine) {
	anomalyRoutes := r.Group("/api/model/anomaly")
	{
		anomalyRoutes.POST("/create", middleware.RequireAuth, controllers.AddAnomalyModel)
		anomalyRoutes.PUT("/update/:id", middleware.RequireAuth, controllers.UpdateAnomalyModel)
		anomalyRoutes.GET("/:id", middleware.RequireAuth, controllers.GetAnomalyModel)
		anomalyRoutes.GET("/device/:id", middleware.RequireAuth, controllers.GetAnomalyModels)
		anomalyRoutes.DELETE("/:id", middleware.RequireAuth, controllers.DeleteAnomalyModel)
	}
}
