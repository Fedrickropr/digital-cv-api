package models

import (
	"github.com/google/uuid"
)

type JwtClaim struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey"`
	TokenID uuid.UUID `gorm:"type:uuid;foreignKey:JwtTokenModel"`
	Claim   string
	Value   string
}
