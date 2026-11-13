package models

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name        string `gorm:"size:100;not null;uniqueIndex" json:"name"`
	Description string `gorm:"size:255" json:"description"`
	IsSystem    bool   `gorm:"not null;default:false" json:"isSystem"`

	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions"`
	Users       []User       `gorm:"many2many:user_roles;" json:"-"`
}
