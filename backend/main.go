package main

import (
	"backend/initializers"
	"backend/services"
	"backend/sockets"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectDB()
}

func main() {

	go sockets.StartTelemetryServer()

	r := gin.Default()

	services.RegisterUserRoutes(r)
	services.RegisterDeviceRoutes(r)
	services.RegisterAutomationRoutes(r)
	services.RegisterAlertRoutes(r)

	r.Run()
}
