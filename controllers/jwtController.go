package controllers

import (
	"digital-cv-api/initializers"
	"digital-cv-api/models"
	"errors"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func GenerateJwt(c *gin.Context) {
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

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"session_uuid": sessionUuid,
		"iat":          jwt.NewNumericDate(time.Now()),
		"exp":          jwt.NewNumericDate(time.Now().Add(60 * 24 * time.Hour)),
	})

	tokenString, err := jwtToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error": "Could not generate token"})
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
