package main

import (
	"backend/initializers"
	"backend/middleware"
	"backend/pipelines"
	"backend/services"
	"backend/sockets"
	"embed"
	"io"
	"io/fs"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

//go:embed dist/*
var frontendEmbed embed.FS

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectDB()
}

func serveIndex(c *gin.Context, dist fs.FS) {
	f, err := dist.Open("index.html")
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Data(http.StatusOK, "text/html; charset=utf-8", data)
}

func main() {

	go sockets.StartMQTTBroker()
	go sockets.StartStreamingServer()
	go sockets.StartAlertsHandler()
	go sockets.StartAutomationsHandler()
	time.Sleep(1 * time.Second)
	go sockets.StartCmdClient()

	go pipelines.StartAnomalyDetectionPipeline()

	r := gin.Default()
	r.RedirectTrailingSlash = false
	r.RedirectFixedPath = false
	r.Use(middleware.CORSMiddleware())

	services.RegisterUserRoutes(r)
	services.RegisterDeviceRoutes(r)
	services.RegisterAutomationRoutes(r)
	services.RegisterAlertRoutes(r)
	services.RegisterDashboardRoutes(r)
	services.RegisterWidgetRoutes(r)
	services.RegisterAnomalyModelRoutes(r)
	services.RegisterRoutineRoutes(r)

	dist, err := fs.Sub(frontendEmbed, "dist")
	if err != nil {
		panic("Failed to read embedded dist folder: " + err.Error())
	}

	fileServer := http.FileServer(http.FS(dist))

	r.GET("/assets/*filepath", func(c *gin.Context) {
		fileServer.ServeHTTP(c.Writer, c.Request)
	})

	r.GET("/", func(c *gin.Context) {
		serveIndex(c, dist)
	})

	r.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.JSON(404, gin.H{"error": "API route not found"})
			return
		}
		serveIndex(c, dist)
	})

	r.Run()
}
