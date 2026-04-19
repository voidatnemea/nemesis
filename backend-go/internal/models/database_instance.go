package models

import "time"

type DatabaseInstance struct {
	ID                uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name              string    `gorm:"type:varchar(255)" json:"name"`
	NodeID            int       `gorm:"index" json:"node_id"`
	DatabaseType      string    `gorm:"type:varchar(50);default:'mysql'" json:"database_type"`
	DatabasePort      int       `gorm:"default:3306" json:"database_port"`
	DatabaseUsername  string    `gorm:"type:varchar(255)" json:"database_username"`
	DatabasePassword  string    `gorm:"type:varchar(255)" json:"-"`
	DatabaseHost      string    `gorm:"type:varchar(255)" json:"database_host"`
	DatabaseSubdomain string    `gorm:"type:varchar(255)" json:"database_subdomain"`
	MaxDatabases      int       `gorm:"default:0" json:"max_databases"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (DatabaseInstance) TableName() string { return "featherpanel_databases" }
