package services

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutineRoutes(r *gin.Engine) {
	RoutineRoutes := r.Group("/api/routine")
	{
		RoutineRoutes.POST("/create", middleware.RequireAuth, controllers.AddRoutine)
		RoutineRoutes.GET("/:id", middleware.RequireAuth, controllers.GetRoutine)
		RoutineRoutes.GET("/", middleware.RequireAuth, controllers.GetRoutines)
		RoutineRoutes.DELETE("/delete/:id", middleware.RequireAuth, controllers.DeleteRoutine)
	}
}
