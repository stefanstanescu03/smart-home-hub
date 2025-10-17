package main

import (
	"backend/initializers"
	"backend/services"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectDB()
}

func main() {
	r := gin.Default()

	services.RegisterUserRoutes(r)

	r.Run()
}
