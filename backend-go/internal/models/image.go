package models

import "time"

type Image struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(255)" json:"name"`
	URL       string    `gorm:"type:varchar(500)" json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Image) TableName() string { return "featherpanel_images" }
