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

	// TODO generic session middleware?
	// TODO kind of pointless, just a few calls and no need for security
	var sessionUuid uuid.UUID
	// Create a UUID for this user's session in a cookie if they did not have one already
	if err != nil {
		sessionUuid = uuid.New()
		initializers.DB.Create(&models.Session{ID: sessionUuid})

		c.Header("x-new-session", "true")
	}

	if sessionUuid == uuid.Nil {
		sessionUuid, err = extractUuidFromToken(cookie, c)
		c.SetCookie("jwt", "", 1, "/", "", false, true)

		// Also create a new session if theirs is not valid
		if err != nil {
			log.Println(err)
			sessionUuid = uuid.New()
			initializers.DB.Create(&models.Session{ID: sessionUuid})
		}
	}

	tokenString, shouldReturn := services.GenerateJWT(sessionUuid, false, c)
	initializers.SyncDatabase()
	if shouldReturn {
		return
	}

	// Optional
	name := c.Query("name")

	// Store the token in the database
	jwtTokenObj := models.JwtToken{
		ID:          uuid.New(),
		Name:        name,
		SessionUuid: sessionUuid,
		Token:       tokenString,
	}
	initializers.DB.Create(&jwtTokenObj)

	// Set the session UUID cookie
	c.SetCookie("jwt", tokenString, 60*60*24*60, "/", "", false, true)

	// For easy updates
	c.Header("Authorization", "Bearer "+tokenString)

	c.Status(200)
}

func GetJwts(c *gin.Context) {
	cookie, err := c.Cookie("jwt")
	if err != nil {
		c.JSON(400, gin.H{"error": "Create a session first"})
		return
	}

	sessionUuid, err := extractUuidFromToken(cookie, c)
	if err != nil {
		c.JSON(400, gin.H{"error": "Create a session first"})
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
	cookie, err := c.Cookie("jwt")
	if err != nil {
		c.JSON(400, gin.H{"error": "Create a session first"})
		return
	}

	sessionUuid, err := extractUuidFromToken(cookie, c)
	if err != nil {
		c.JSON(400, gin.H{"error": "Create a session first"})
		return
	}

	var jwtToken models.JwtToken
	if err := initializers.DB.Where("id = ? AND session_uuid = ?", c.Param("id"), sessionUuid).First(&jwtToken).Error; err != nil {
		c.JSON(404, gin.H{"error": "JWT not found"})
		return
	}

	// create new jwt token

	db := initializers.DB
	db.Model(&models.JwtToken{}).Where("session_uuid = ?", sessionUuid).Update("active", false)

	newJwtTokenString, shouldReturn := services.GenerateJWT(sessionUuid, true, c)
	if shouldReturn {
		c.JSON(500, gin.H{"error": "Could not generate new JWT"})
		return
	}

	if err := initializers.DB.Save(&jwtToken).Error; err != nil {
		c.JSON(500, gin.H{"error": "Could not update JWT"})
		return
	}

	c.JSON(200, gin.H{"token": newJwtTokenString, "active": jwtToken.Active})

}

func extractUuidFromToken(cookie string, c *gin.Context) (uuid.UUID, error) {
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
