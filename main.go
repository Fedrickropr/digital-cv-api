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
	jwt.GET("", controllers.GetJwts)
	jwt.POST("", controllers.CreateJwt)
	jwt.PUT("/:id", controllers.UpdateJwt)
	jwt.GET("/:id/claims", controllers.GetJwtClaimsById)
	jwt.DELETE("/:id", controllers.DeleteJwt)

	claim := jwt.Group("claim")
	claim.GET("", controllers.GetJwtClaims)
	claim.POST("", controllers.AddJwtClaim)
	claim.PUT("/:id", controllers.EditJwtClaim)
	claim.DELETE("/:id", controllers.DeleteJwtClaim)

	router.Run("localhost:3000")
}
