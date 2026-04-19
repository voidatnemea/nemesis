package models

import "time"

type Node struct {
	ID               uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID             string    `gorm:"uniqueIndex;type:varchar(36)" json:"uuid"`
	Name             string    `gorm:"type:varchar(255)" json:"name"`
	Description      string    `gorm:"type:text" json:"description"`
	LocationID       int       `gorm:"index" json:"location_id"`
	FQDN             string    `gorm:"type:varchar(255)" json:"fqdn"`
	Public           string    `gorm:"type:enum('true','false');default:'true'" json:"public"`
	Scheme           string    `gorm:"type:enum('http','https');default:'https'" json:"scheme"`
	BehindProxy      string    `gorm:"type:enum('true','false');default:'false'" json:"behind_proxy"`
	Memory           int64     `json:"memory"`
	MemoryOvercommit int       `gorm:"default:0" json:"memory_overcommit"`
	Disk             int64     `json:"disk"`
	DiskOvercommit   int       `gorm:"default:0" json:"disk_overcommit"`
	DaemonBase       string    `gorm:"type:varchar(255);default:'/var/lib/pterodactyl/volumes'" json:"daemon_base"`
	DaemonSFTP       int       `gorm:"default:2022" json:"daemon_sftp"`
	DaemonListen     int       `gorm:"default:8080" json:"daemon_listen"`
	DaemonToken      string    `gorm:"type:varchar(255)" json:"-"`
	MaintenanceMode  string    `gorm:"type:enum('true','false');default:'false'" json:"maintenance_mode"`
	UploadSize       int       `gorm:"default:100" json:"upload_size"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (Node) TableName() string { return "featherpanel_nodes" }
