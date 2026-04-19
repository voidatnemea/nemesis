package models

import "time"

type InstalledPlugin struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"uniqueIndex;type:varchar(255)" json:"name"`
	Version     string    `gorm:"type:varchar(50)" json:"version"`
	Author      string    `gorm:"type:varchar(255)" json:"author"`
	Description string    `gorm:"type:text" json:"description"`
	Enabled     string    `gorm:"type:enum('true','false');default:'true'" json:"enabled"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (InstalledPlugin) TableName() string { return "featherpanel_installed_plugins" }
