package main

import (
	"backend/initializers"
	"backend/middleware"
	"backend/pipelines"
	"backend/services"
	"backend/sockets"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectDB()
}

func main() {

	// go sockets.StartTelemetryServer()
	go sockets.StartMQTTBroker()
	go sockets.StartStreamingServer()
	go sockets.StartAlertsHandler()
	go sockets.StartAutomationsHandler()
	time.Sleep(1 * time.Second)
	go sockets.StartCmdClient()

	go pipelines.StartAnomalyDetectionPipeline()

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	services.RegisterUserRoutes(r)
	services.RegisterDeviceRoutes(r)
	services.RegisterAutomationRoutes(r)
	services.RegisterAlertRoutes(r)
	services.RegisterDashboardRoutes(r)
	services.RegisterWidgetRoutes(r)
	services.RegisterAnomalyModelRoutes(r)

	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "if you see this the backend connection works",
		})
	})

	r.Run()
}
