package models

import (
	"time"

	"gorm.io/gorm"
)

type UserSshKey struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      int            `gorm:"index" json:"user_id"`
	Name        string         `gorm:"type:varchar(255)" json:"name"`
	PublicKey   string         `gorm:"type:text" json:"public_key"`
	Fingerprint string         `gorm:"type:varchar(255)" json:"fingerprint"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (UserSshKey) TableName() string { return "featherpanel_user_ssh_keys" }

type UserPreference struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserUUID  string    `gorm:"type:varchar(36);index" json:"user_uuid"`
	Key       string    `gorm:"type:varchar(255)" json:"key"`
	Value     string    `gorm:"type:text" json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (UserPreference) TableName() string { return "featherpanel_user_preferences" }
