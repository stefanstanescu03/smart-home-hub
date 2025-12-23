package services

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterWidgetRoutes(r *gin.Engine) {
	widgetRoutes := r.Group("/widget")
	{
		widgetRoutes.POST("/create", middleware.RequireAuth, controllers.AddWidget)
		widgetRoutes.GET("/:id", middleware.RequireAuth, controllers.GetWidgets)
		widgetRoutes.DELETE("/delete/:id", middleware.RequireAuth, controllers.DeleteWidget)
		widgetRoutes.GET("/public/:id", controllers.GetPublicWidgets)
	}
}
