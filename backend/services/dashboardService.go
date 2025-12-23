package services

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterDashboardRoutes(r *gin.Engine) {
	dashboardRoutes := r.Group("/dashboard")
	{
		dashboardRoutes.POST("/create", middleware.RequireAuth, controllers.AddDashboard)
		dashboardRoutes.GET("/:id", middleware.RequireAuth, controllers.GetDashboard)
		dashboardRoutes.GET("/", middleware.RequireAuth, controllers.GetDashboards)
		dashboardRoutes.DELETE("/delete/:id", middleware.RequireAuth, controllers.DeleteDashboard)
		dashboardRoutes.PUT("/update/:id", middleware.RequireAuth, controllers.UpdateDashboard)
		dashboardRoutes.GET("/public", controllers.GetPublicDashboards)
		dashboardRoutes.GET("/public/:id", controllers.GetDashboardPublic)
	}
}
