package models

import "time"

type Mount struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"type:varchar(255)" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Source      string    `gorm:"type:varchar(255)" json:"source"`
	Target      string    `gorm:"type:varchar(255)" json:"target"`
	ReadOnly    string    `gorm:"type:enum('true','false');default:'false'" json:"read_only"`
	UserMountable string  `gorm:"type:enum('true','false');default:'false'" json:"user_mountable"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Mount) TableName() string { return "featherpanel_mounts" }
