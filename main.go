package main

import (
	"digital-cv-api/controllers"
	"digital-cv-api/initializers"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	initializers.SyncDatabase()
}

func main() {
	router := gin.Default()
	api := router.Group("api")
	api.POST("/jwt", controllers.GenerateJwt)

	router.Run("localhost:3000")
}
