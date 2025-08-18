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
	jwt := api.Group("jwt")
	claim := jwt.Group("claim")

	jwt.GET("", controllers.GetJwts)
	jwt.POST("", controllers.CreateJwt)
	jwt.PUT("/:id", controllers.UpdateJwt)
	claim.POST("", controllers.AddJwtClaim)

	router.Run("localhost:3000")
}
