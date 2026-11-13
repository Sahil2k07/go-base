package models

import "gorm.io/gorm"

type Permission struct {
	gorm.Model
	Name        string `gorm:"size:100;not null;uniqueIndex" json:"name"`
	Description string `gorm:"size:255" json:"description"`
	Roles       []Role `gorm:"many2many:role_permissions;" json:"-"`
}
