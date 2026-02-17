package services

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine) {
	userRoutes := r.Group("/api/user")
	{
		userRoutes.POST("/signup", controllers.Signup)
		userRoutes.POST("/login", controllers.Login)
		userRoutes.GET("/info", middleware.RequireAuth, controllers.ShowUser)
		userRoutes.PUT("/update", middleware.RequireAuth, controllers.ModifyUser)
		userRoutes.DELETE("/delete/:id", middleware.RequireAuth, controllers.DeleteUser)
		userRoutes.GET("/all", middleware.RequireAuth, controllers.GetUsers)
		userRoutes.PUT("/toggle-role/:id", middleware.RequireAuth, controllers.ToggleRole)
	}
}
