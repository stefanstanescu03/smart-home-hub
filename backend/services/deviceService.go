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
		deviceRoutes.PUT("/update/:id", middleware.RequireAuth, controllers.UpdateDevice)
		deviceRoutes.GET("/:id", middleware.RequireAuth, controllers.GetDevice)
		deviceRoutes.GET("/", middleware.RequireAuth, controllers.GetDevices)
		deviceRoutes.DELETE("/delete/:id", middleware.RequireAuth, controllers.DeleteDevice)
		deviceRoutes.GET("/ping/:ident", controllers.IsDeviceConnected)
		deviceRoutes.GET("/public", controllers.GetPublicDevices)
	}
}
