package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID                uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Username          string         `gorm:"uniqueIndex;type:varchar(255)" json:"username"`
	FirstName         string         `gorm:"type:varchar(255)" json:"first_name"`
	LastName          string         `gorm:"type:varchar(255)" json:"last_name"`
	Email             string         `gorm:"uniqueIndex;type:varchar(255)" json:"email"`
	Password          string         `gorm:"type:varchar(255)" json:"-"`
	UUID              string         `gorm:"uniqueIndex;type:varchar(36)" json:"uuid"`
	RememberToken      string         `gorm:"type:varchar(255)" json:"-"`
	RoleID            int            `gorm:"default:1" json:"role_id"`
	Avatar            string         `gorm:"type:varchar(255)" json:"avatar"`
	FirstIP           string         `gorm:"type:varchar(45)" json:"first_ip"`
	LastIP            string         `gorm:"type:varchar(45)" json:"last_ip"`
	Banned            string         `gorm:"type:enum('true','false');default:'false'" json:"banned"`
	TwoFAEnabled      string         `gorm:"type:enum('true','false');default:'false'" json:"two_fa_enabled"`
	TwoFAKey          string         `gorm:"type:varchar(255)" json:"-"`
	ExternalID        string         `gorm:"type:varchar(255)" json:"external_id"`
	Deleted           string         `gorm:"type:enum('true','false');default:'false'" json:"deleted"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "featherpanel_users"
}
