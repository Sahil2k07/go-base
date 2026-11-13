package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserSession struct {
	gorm.Model
	UserID            uint      `gorm:"not null;index" json:"userId"`
	SessionIdentifier uuid.UUID `gorm:"type:uuid;not null;uniqueIndex" json:"sessionIdentifier"`
	RefreshTokenHash  string    `gorm:"not null" json:"-"`
	ExpiresAt         time.Time `gorm:"not null;index" json:"expiresAt"`

	User User `gorm:"foreignKey:UserID" json:"-"` // belongs-to backlink
}
