package models

import "time"

type Backup struct {
	ID           uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	ServerID     int        `gorm:"index" json:"server_id"`
	UUID         string     `gorm:"uniqueIndex;type:varchar(36)" json:"uuid"`
	Name         string     `gorm:"type:varchar(255)" json:"name"`
	IgnoredFiles string     `gorm:"type:text" json:"ignored_files"`
	Disk         string     `gorm:"type:varchar(50)" json:"disk"`
	SHA256Hash   string     `gorm:"type:varchar(64)" json:"sha256_hash"`
	Bytes        int64      `gorm:"default:0" json:"bytes"`
	IsSuccessful *bool      `json:"is_successful"`
	IsLocked     string     `gorm:"type:enum('true','false');default:'false'" json:"is_locked"`
	CompletedAt  *time.Time `json:"completed_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

func (Backup) TableName() string { return "featherpanel_server_backups" }
