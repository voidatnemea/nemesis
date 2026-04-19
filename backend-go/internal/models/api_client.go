package models

import (
	"time"
)

type ApiClient struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserUUID    string    `gorm:"type:varchar(36);index" json:"user_uuid"`
	Name        string    `gorm:"type:varchar(255)" json:"name"`
	PublicKey   string    `gorm:"uniqueIndex;type:varchar(255)" json:"public_key"`
	PrivateKey  string    `gorm:"uniqueIndex;type:varchar(255)" json:"-"`
	AllowedIps  string    `gorm:"type:text" json:"allowed_ips"`
	LastUsedAt  *time.Time `json:"last_used_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (ApiClient) TableName() string {
	return "featherpanel_api_clients"
}
