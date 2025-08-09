package controllers

import (
	"digital-cv-api/initializers"
	"digital-cv-api/models"
	"errors"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func GenerateJwt(c *gin.Context) {
	cookie, err := c.Cookie("jwt")

	var sessionUuid uuid.UUID
	// Create a UUID for this user's session in a cookie if they did not have one already
	if err != nil {
		sessionUuid = uuid.New()
		initializers.DB.Create(&models.Session{ID: sessionUuid})

		c.Header("x-new-session", "true")
	}

	if sessionUuid == uuid.Nil {
		sessionUuid, err = extractUuidFromToken(cookie, sessionUuid, c)

		if err != nil {
			log.Println(err)
			c.JSON(500, gin.H{"error": "Found token but no valid UUID"})
			return
		}
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"session_uuid": sessionUuid,
	})

	tokenString, err := jwtToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error": "Could not generate token"})
		return
	}

	// TODO store the token to the db
	// TODO maybe now implement a service layer?

	// Set the session UUID cookie
	c.SetCookie("jwt", tokenString, 60*60*24*60, "/", "", false, true)

	// For easy updates
	c.Header("Authorization", "Bearer "+tokenString)

	c.Status(200)
}

func extractUuidFromToken(cookie string, sessionUuid uuid.UUID, c *gin.Context) (uuid.UUID, error) {
	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	}, jwt.WithValidMethods([]string{"HS256"}))

	if err != nil {
		log.Println(err)
		return uuid.Nil, errors.New("could not extract UUID from token")
	}

	if err == nil && token.Valid {
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if sid, ok := claims["session_uuid"].(string); ok {
				sessionUuid, err = uuid.Parse(sid)
				if err != nil {
					// invalid UUID string in claim, fallback to new UUID
					sessionUuid = uuid.New()
					initializers.DB.Create(&models.Session{ID: sessionUuid})
					c.Header("x-new-session", "true")
				}
			}
		}
	}
	return sessionUuid, nil
}
