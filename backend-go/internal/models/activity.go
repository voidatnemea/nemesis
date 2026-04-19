package models

import "time"

type Activity struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserUUID    string    `gorm:"type:varchar(36);index" json:"user_uuid"`
	Event       string    `gorm:"type:varchar(255)" json:"event"`
	Description string    `gorm:"type:text" json:"description"`
	IP          string    `gorm:"type:varchar(45)" json:"ip"`
	Metadata    string    `gorm:"type:json" json:"metadata"`
	CreatedAt   time.Time `json:"created_at"`
}

func (Activity) TableName() string { return "featherpanel_activity" }

type TimedTask struct {
	ID          uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string     `gorm:"type:varchar(255)" json:"name"`
	Expression  string     `gorm:"type:varchar(255)" json:"expression"`
	IsActive    string     `gorm:"type:enum('true','false');default:'true'" json:"is_active"`
	LastRunAt   *time.Time `json:"last_run_at"`
	NextRunAt   *time.Time `json:"next_run_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (TimedTask) TableName() string { return "featherpanel_timed_tasks" }
