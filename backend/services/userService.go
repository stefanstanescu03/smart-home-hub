package services

import (
	"backend/controllers"
	"backend/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine) {
	userRoutes := r.Group("/user")
	{
		userRoutes.POST("/signup", controllers.Signup)
		userRoutes.POST("/login", controllers.Login)
		userRoutes.GET("/info", middleware.RequireAuth, controllers.ShowUser)
	}
}
