package models

import "time"

type Location struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"type:varchar(255)" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	FlagCode    string    `gorm:"type:varchar(10)" json:"flag_code"`
	Type        string    `gorm:"type:enum('game','vps','web');default:'game'" json:"type"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Location) TableName() string { return "featherpanel_locations" }
