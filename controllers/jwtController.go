package controllers

import (
	"digital-cv-api/initializers"
	"digital-cv-api/models"
	"digital-cv-api/services"
	"errors"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func CreateJwt(c *gin.Context) {
	cookie, err := c.Cookie("jwt")
	var sessionUuid uuid.UUID

	// Create a UUID for this user's session in a cookie if they did not have one already
	if err != nil {
		sessionUuid = uuid.New()
		initializers.DB.Create(&models.Session{ID: sessionUuid})

		c.Header("x-new-session", "true")
	}

	// If we didnt create one, i.e. cookie existed, parse the cookie for it
	if sessionUuid == uuid.Nil {
		sessionUuid, err = extractUuidFromToken(cookie)
		c.SetCookie("jwt", "", 1, "/", "", false, true)

		// Insert it if we didnt have it stored
		if err != nil {
			sessionUuid = uuid.New()
			initializers.DB.Create(&models.Session{ID: sessionUuid})
		}
	}

	// Insert token description to DB
	name := c.Query("name")
	jwtTokenObj := models.JwtToken{
		ID:          uuid.New(),
		Name:        name,
		SessionUuid: sessionUuid,
		Active:      true,
	}
	initializers.DB.Create(&jwtTokenObj)

	// Generate the token
	tokenString, err := services.GenerateJWT(sessionUuid, jwtTokenObj.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	// Set the session UUID cookie
	c.SetCookie("jwt", tokenString, 60*60*24*60, "/", "", false, true)

	// Allows frontend to always access active token
	c.Header("Authorization", "Bearer "+tokenString)

	c.Status(200)
}

func GetJwts(c *gin.Context) {
	sessionUuid, err := getSessionUuid(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var jwtList []models.JwtToken
	res := initializers.DB.Where("session_uuid = ?", sessionUuid).Find(&jwtList)

	if res.Error != nil {
		c.JSON(500, gin.H{"error": "Could not get JWTs"})
		return
	}

	c.JSON(200, jwtList)
}

func UpdateJwt(c *gin.Context) {
	sessionUuid, err := getSessionUuid(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var jwtToken models.JwtToken
	if err := initializers.DB.Where("id = ? AND session_uuid = ?", c.Param("id"), sessionUuid).First(&jwtToken).Error; err != nil {
		c.JSON(404, gin.H{"error": "JWT not found"})
		return
	}

	// Set all tokens inactive
	initializers.DB.Model(&models.JwtToken{}).Where("session_uuid = ?", sessionUuid).Update("active", false)

	// Set the to update token to active
	jwtToken.Active = true
	if err = initializers.DB.Save(&jwtToken).Error; err != nil {
		c.JSON(500, gin.H{"error": "Could not update JWT"})
		return
	}

	newJwtTokenString, err := services.GenerateJWT(sessionUuid, jwtToken.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"token": newJwtTokenString, "active": jwtToken.Active})
}

func GetJwtContents(c *gin.Context) {
	sessionUuid, err := getSessionUuid(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	jwtUuid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid UUID format " + err.Error()})
		return
	}

	tokenString, err := services.GenerateJWT(sessionUuid, jwtUuid)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
	}

	c.JSON(200, gin.H{"token": tokenString})
}

func getSessionUuid(c *gin.Context) (uuid.UUID, error) {
	cookie, err := c.Cookie("jwt")
	if err != nil {
		return uuid.UUID{}, errors.New("no cookie found")
	}

	sessionUuid, err := extractUuidFromToken(cookie)
	if err != nil {
		return uuid.UUID{}, errors.New("create a session first")
	}
	return sessionUuid, nil
}

func extractUuidFromToken(cookie string) (uuid.UUID, error) {
	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	}, jwt.WithValidMethods([]string{"HS256"}))

	if err != nil {
		log.Println(err)
		return handleUuidError(err, "could not parse token")
	}

	if !token.Valid {
		log.Println(err)
		return handleUuidError(err, "token was not valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return handleUuidError(err, "invalid token claims")
	}

	sessionId, ok := claims["session_uuid"].(string)
	if !ok {
		log.Println(err)
		return handleUuidError(err, "no session uuid found")
	}

	sessionUuid, err := uuid.Parse(sessionId)
	if err != nil {
		log.Println(err)
		return uuid.Nil, errors.New("could not parse UUID from token")
	}

	return sessionUuid, nil
}

func handleUuidError(err error, msg string) (uuid.UUID, error) {
	log.Println(err)
	return uuid.Nil, errors.New(msg)
}
