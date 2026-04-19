package models

import "time"

type ZeroTrustHash struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Hash        string    `gorm:"type:varchar(64);uniqueIndex" json:"hash"`
	FileName    string    `gorm:"type:varchar(255)" json:"file_name"`
	Description string    `gorm:"type:text" json:"description"`
	Status      string    `gorm:"type:enum('pending','confirmed');default:'pending'" json:"status"`
	DetectedAt  time.Time `json:"detected_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (ZeroTrustHash) TableName() string { return "featherpanel_zerotrust_hashes" }

type ZeroTrustScanExecution struct {
	ID              uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Status          string     `gorm:"type:enum('running','completed','failed');default:'running'" json:"status"`
	ServersScanned  int        `json:"servers_scanned"`
	DetectionsFound int        `json:"detections_found"`
	DurationMs      int64      `json:"duration_ms"`
	StartedAt       time.Time  `json:"started_at"`
	FinishedAt      *time.Time `json:"finished_at"`
	CreatedAt       time.Time  `json:"created_at"`
}

func (ZeroTrustScanExecution) TableName() string { return "featherpanel_zerotrust_executions" }

type ZeroTrustScanLog struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ExecutionID   uint      `gorm:"index" json:"execution_id"`
	ServerUUID    string    `gorm:"type:varchar(36);index" json:"server_uuid"`
	ServerName    string    `gorm:"type:varchar(255)" json:"server_name"`
	Status        string    `gorm:"type:enum('clean','suspicious','error');default:'clean'" json:"status"`
	FilesScanned  int       `json:"files_scanned"`
	Detections    string    `gorm:"type:json" json:"detections"`
	ErrorMessage  string    `gorm:"type:text" json:"error_message"`
	CreatedAt     time.Time `json:"created_at"`
}

func (ZeroTrustScanLog) TableName() string { return "featherpanel_zerotrust_scan_logs" }
