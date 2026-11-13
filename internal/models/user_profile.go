package models

import "gorm.io/gorm"

type UserProfile struct {
	gorm.Model
	UserID    uint   `gorm:"not null;uniqueIndex" json:"userId"`
	FirstName string `gorm:"size:100;not null" json:"firstName"`
	LastName  string `gorm:"size:100;not null" json:"lastName"`
	Phone     string `gorm:"size:20" json:"phone"`
	Gender    string `gorm:"size:20" json:"gender"`

	User *User `gorm:"foreignKey:UserID" json:"-"`
}
