package controllers

import (
	"digital-cv-api/initializers"
	"digital-cv-api/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetJwtClaims(c *gin.Context) {
	var jwtClaims []models.JwtClaim
	if err := initializers.DB.Find(&jwtClaims).Error; err != nil {
		c.JSON(500, gin.H{"error": "Could not retrieve JWT claims"})
		return
	}

	c.JSON(200, jwtClaims)
}

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

	c.JSON(201, jwtClaim)

}

func EditJwtClaim(c *gin.Context) {

	var jwtClaim models.JwtClaim
	if err := c.ShouldBindJSON(&jwtClaim); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	jwtClaimID := c.Param("id")
	if err := initializers.DB.Model(&models.JwtClaim{}).Where("id = ?", jwtClaimID).Updates(jwtClaim).Error; err != nil {
		c.JSON(500, gin.H{"error": "Could not update JWT claim"})
		return
	}

	c.JSON(200, gin.H{"message": "JWT claim updated successfully"})
}

func DeleteJwtClaim(c *gin.Context) {
	jwtClaimID := c.Param("id")

	if err := initializers.DB.Where("id = ?", jwtClaimID).Delete(&models.JwtClaim{}).Error; err != nil {
		c.JSON(500, gin.H{"error": "Could not delete JWT claim"})
		return
	}

	c.JSON(200, gin.H{"message": "JWT claim deleted successfully"})
}
