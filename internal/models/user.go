package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email        string `gorm:"size:150;not null;uniqueIndex" json:"email"`
	PasswordHash string `gorm:"not null" json:"-"`

	Profile  UserProfile   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"profile"`
	Sessions []UserSession `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Roles    []Role        `gorm:"many2many:user_roles;" json:"roles"`
}
