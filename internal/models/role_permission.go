package models

import "time"

type RolePermission struct {
	RoleID       uint      `gorm:"primaryKey" json:"roleId"`
	PermissionID uint      `gorm:"primaryKey;index" json:"permissionId"`
	CreatedAt    time.Time `json:"createdAt"`

	Role       Role       `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE" json:"-"`
	Permission Permission `gorm:"foreignKey:PermissionID;constraint:OnDelete:CASCADE" json:"-"`
}
