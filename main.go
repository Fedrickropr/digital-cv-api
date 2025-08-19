package main

import (
	"digital-cv-api/controllers"
	"digital-cv-api/initializers"
	"digital-cv-api/middleware"

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
	api.Handlers = append(api.Handlers, middleware.HandleCors())

	jwt := api.Group("jwt")
	jwt.GET("", controllers.GetJwts)
	jwt.POST("", controllers.CreateJwt)
	jwt.PUT("/:id", controllers.UpdateJwt)
	jwt.GET("/:id", controllers.GetJwtContents)

	claim := jwt.Group("claim")
	claim.POST("", controllers.AddJwtClaim)

	router.Run("localhost:3000")
}
