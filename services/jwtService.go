package services

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func GenerateJWT(sessionUuid uuid.UUID, active bool, c *gin.Context) (string, bool) {
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"session_uuid": sessionUuid,
		"active":       active,
		"iat":          jwt.NewNumericDate(time.Now()),
		"exp":          jwt.NewNumericDate(time.Now().Add(60 * 24 * time.Hour)),
	})

	tokenString, err := jwtToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error": "Could not generate token"})
		return "", true
	}
	return tokenString, false
}
