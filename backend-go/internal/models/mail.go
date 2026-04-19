package models

import "time"

type MailTemplate struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"uniqueIndex;type:varchar(255)" json:"name"`
	Subject   string    `gorm:"type:varchar(255)" json:"subject"`
	Body      string    `gorm:"type:longtext" json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (MailTemplate) TableName() string { return "featherpanel_mail_templates" }

type MailQueue struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserUUID  string    `gorm:"type:varchar(36);index" json:"user_uuid"`
	Subject   string    `gorm:"type:varchar(255)" json:"subject"`
	Body      string    `gorm:"type:longtext" json:"body"`
	Status    string    `gorm:"type:enum('pending','sent','failed');default:'pending'" json:"status"`
	SentAt    *time.Time `json:"sent_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (MailQueue) TableName() string { return "featherpanel_mail_queue" }
