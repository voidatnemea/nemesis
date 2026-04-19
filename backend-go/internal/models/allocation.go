package models

import "time"

type Allocation struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	NodeID    int       `gorm:"index" json:"node_id"`
	ServerID  *int      `gorm:"index" json:"server_id"`
	IP        string    `gorm:"type:varchar(45)" json:"ip"`
	Port      int       `json:"port"`
	Alias     string    `gorm:"type:varchar(255)" json:"alias"`
	Notes     string    `gorm:"type:text" json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Allocation) TableName() string { return "featherpanel_allocations" }
