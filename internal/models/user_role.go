package models

import "time"

type UserRole struct {
	UserID    uint      `gorm:"primaryKey" json:"userId"`
	RoleID    uint      `gorm:"primaryKey;index" json:"roleId"`
	CreatedAt time.Time `json:"createdAt"`

	User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"-"`
	Role Role `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE" json:"-"`
}
