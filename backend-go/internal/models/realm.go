package models

import "time"

type Realm struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"type:varchar(255)" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Realm) TableName() string { return "featherpanel_realms" }
