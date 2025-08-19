package models

import (
	"github.com/google/uuid"
)

type JwtToken struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name        string
	SessionUuid uuid.UUID `gorm:"type:uuid;foreignKey:SessionModel"`
	Active      bool      `gorm:"type:boolean"`
}
