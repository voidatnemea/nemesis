package models

import "time"

type Notification struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Title           string    `gorm:"type:varchar(255)" json:"title"`
	MessageMarkdown string    `gorm:"type:text" json:"message_markdown"`
	Type            string    `gorm:"type:enum('info','success','warning','error');default:'info'" json:"type"`
	IsDismissible   string    `gorm:"type:enum('true','false');default:'true'" json:"is_dismissible"`
	IsSticky        string    `gorm:"type:enum('true','false');default:'false'" json:"is_sticky"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (Notification) TableName() string { return "featherpanel_notifications" }
