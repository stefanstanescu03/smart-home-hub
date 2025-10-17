package main

import (
	"backend/initializers"
	"backend/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectDB()
}

func main() {
	initializers.DB.AutoMigrate(&models.User{})
	initializers.DB.AutoMigrate(&models.Device{})
	initializers.DB.AutoMigrate(&models.Automation{})
	initializers.DB.AutoMigrate(&models.Alert{})
}
