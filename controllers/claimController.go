package controllers

import (
	"digital-cv-api/initializers"
	"digital-cv-api/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AddJwtClaim(c *gin.Context) {

	var jwtClaim models.JwtClaim
	if err := c.ShouldBindJSON(&jwtClaim); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	jwtClaim.ID = uuid.New()
	if err := initializers.DB.Create(&jwtClaim).Error; err != nil {
		c.JSON(500, gin.H{"error": "Could not create JWT claim"})
		return
	}

}
