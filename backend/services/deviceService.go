package services

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterDeviceRoutes(r *gin.Engine) {
	deviceRoutes := r.Group("/device")
	{
		deviceRoutes.POST("/create", middleware.RequireAuth, controllers.AddDevice)
		deviceRoutes.PUT("/update/:name", middleware.RequireAuth, controllers.UpdateDevice)
		deviceRoutes.GET("/:name", middleware.RequireAuth, controllers.GetDevice)
		deviceRoutes.GET("/", middleware.RequireAuth, controllers.GetDevices)
		deviceRoutes.DELETE("/delete/:name", middleware.RequireAuth, controllers.DeleteDevice)
	}
}
