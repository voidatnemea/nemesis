package models

import "time"

type Role struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"uniqueIndex;type:varchar(255)" json:"name"`
	DisplayName string    `gorm:"type:varchar(255)" json:"display_name"`
	Color       string    `gorm:"type:varchar(50);default:'#6366f1'" json:"color"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Role) TableName() string { return "featherpanel_roles" }

type Permission struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	RoleID     int       `gorm:"index" json:"role_id"`
	Permission string    `gorm:"type:varchar(255)" json:"permission"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (Permission) TableName() string { return "featherpanel_permissions" }
