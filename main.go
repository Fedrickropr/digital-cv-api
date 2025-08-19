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
	jwt.GET("/:id", controllers.GetJwtContents)
	jwt.GET("/:id/claims", controllers.GetJwtClaimsById)
	jwt.POST("", controllers.CreateJwt)
	jwt.PUT("/:id", controllers.UpdateJwt)
	jwt.DELETE("/:id", controllers.DeleteJwt)

	claim := jwt.Group("claim")
	claim.GET("", controllers.GetJwtClaims)
	claim.POST("", controllers.AddJwtClaim)
	claim.PUT("/:id", controllers.EditJwtClaim)
	claim.DELETE("/:id", controllers.DeleteJwtClaim)

	router.Run("localhost:3000")
}
