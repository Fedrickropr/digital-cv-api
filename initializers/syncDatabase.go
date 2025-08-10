package initializers

import (
	"digital-cv-api/models"
)

func SyncDatabase() {
	DB.AutoMigrate(&models.Session{})
	DB.AutoMigrate(&models.JwtToken{})
}
