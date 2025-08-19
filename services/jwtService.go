package services

import (
	"digital-cv-api/initializers"
	"digital-cv-api/models"
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

func GenerateJWT(sessionUuid uuid.UUID, tokenUuid uuid.UUID) (string, error) {
	// Get the token
	var tokenData models.JwtToken
	res := initializers.DB.Where("session_uuid = ? AND id = ?", sessionUuid, tokenUuid).First(&tokenData)
	if res.Error != nil {
		return "", errors.New("could not find JWT")
	}

	// Get the tokenClaims
	var tokenClaims []models.JwtClaim
	res = initializers.DB.Where("token_id = ?", tokenData.ID).Find(&tokenClaims)
	if res.Error != nil {
		return "", errors.New("could not get claims")
	}

	tokenContents := jwt.MapClaims{
		"session_uuid": sessionUuid,
		"iat":          jwt.NewNumericDate(time.Now()),
		"exp":          jwt.NewNumericDate(time.Now().Add(60 * 24 * time.Hour)),
	}

	for _, c := range tokenClaims {
		tokenContents[c.Claim] = c.Value
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenContents)

	tokenString, err := jwtToken.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		log.Println(err)
		return "", errors.New("could not generate token")
	}
	return tokenString, nil
}
