package models

import (
	"time"

	"gorm.io/gorm"
)

type Server struct {
	ID               uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID             string         `gorm:"uniqueIndex;type:varchar(36)" json:"uuid"`
	UUIDShort        string         `gorm:"uniqueIndex;type:varchar(8)" json:"uuid_short"`
	NodeID           int            `gorm:"index" json:"node_id"`
	Name             string         `gorm:"type:varchar(255)" json:"name"`
	Description      string         `gorm:"type:text" json:"description"`
	OwnerID          int            `gorm:"index" json:"owner_id"`
	Memory           int64          `json:"memory"`
	Swap             int64          `json:"swap"`
	Disk             int64          `json:"disk"`
	IO               int            `json:"io"`
	CPU              int            `json:"cpu"`
	Threads          string         `gorm:"type:varchar(255)" json:"threads"`
	OOMDisabled      string         `gorm:"type:enum('true','false');default:'false'" json:"oom_disabled"`
	AllocationID     int            `json:"allocation_id"`
	RealmsID         int            `json:"realms_id"`
	SpellID          int            `json:"spell_id"`
	Image            string         `gorm:"type:varchar(500)" json:"image"`
	Startup          string         `gorm:"type:text" json:"startup"`
	Status           string         `gorm:"type:enum('installing','install_failed','suspended','restoring_backup','');default:'installing'" json:"status"`
	Suspended        string         `gorm:"type:enum('true','false');default:'false'" json:"suspended"`
	SkipScripts      string         `gorm:"type:enum('true','false');default:'false'" json:"skip_scripts"`
	ExternalID       string         `gorm:"type:varchar(255)" json:"external_id"`
	DatabaseLimit    int            `gorm:"default:0" json:"database_limit"`
	AllocationLimit  int            `gorm:"default:0" json:"allocation_limit"`
	BackupLimit      int            `gorm:"default:0" json:"backup_limit"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Server) TableName() string { return "featherpanel_servers" }

type ServerVariable struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ServerID      int       `gorm:"index" json:"server_id"`
	VariableID    int       `json:"variable_id"`
	VariableValue string    `gorm:"type:text" json:"variable_value"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (ServerVariable) TableName() string { return "featherpanel_server_variables" }

type ServerActivity struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ServerID    int       `gorm:"index" json:"server_id"`
	UserUUID    string    `gorm:"type:varchar(36);index" json:"user_uuid"`
	Event       string    `gorm:"type:varchar(255)" json:"event"`
	Description string    `gorm:"type:text" json:"description"`
	IP          string    `gorm:"type:varchar(45)" json:"ip"`
	Metadata    string    `gorm:"type:json" json:"metadata"`
	CreatedAt   time.Time `json:"created_at"`
}

func (ServerActivity) TableName() string { return "featherpanel_server_activities" }

type ServerDatabase struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ServerID       int       `gorm:"index" json:"server_id"`
	DatabaseHostID int       `json:"database_host_id"`
	Database       string    `gorm:"type:varchar(255)" json:"database"`
	Username       string    `gorm:"type:varchar(255)" json:"username"`
	Password       string    `gorm:"type:varchar(255)" json:"-"`
	Remote         string    `gorm:"type:varchar(45);default:'%'" json:"remote"`
	MaxConnections int       `gorm:"default:0" json:"max_connections"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (ServerDatabase) TableName() string { return "featherpanel_server_databases" }

type ServerSchedule struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ServerID        int       `gorm:"index" json:"server_id"`
	Name            string    `gorm:"type:varchar(255)" json:"name"`
	CronDayOfWeek   string    `gorm:"type:varchar(10)" json:"cron_day_of_week"`
	CronMonth       string    `gorm:"type:varchar(10)" json:"cron_month"`
	CronDayOfMonth  string    `gorm:"type:varchar(10)" json:"cron_day_of_month"`
	CronHour        string    `gorm:"type:varchar(10)" json:"cron_hour"`
	CronMinute      string    `gorm:"type:varchar(10)" json:"cron_minute"`
	IsActive        string    `gorm:"type:enum('true','false');default:'true'" json:"is_active"`
	IsProcessing    string    `gorm:"type:enum('true','false');default:'false'" json:"is_processing"`
	OnlyWhenOnline  string    `gorm:"type:enum('true','false');default:'false'" json:"only_when_online"`
	LastRunAt       *time.Time `json:"last_run_at"`
	NextRunAt       *time.Time `json:"next_run_at"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (ServerSchedule) TableName() string { return "featherpanel_server_schedules" }

type ServerTransfer struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ServerID      int       `gorm:"index" json:"server_id"`
	OldNodeID     int       `json:"old_node_id"`
	NewNodeID     int       `json:"new_node_id"`
	OldAllocation int       `json:"old_allocation"`
	NewAllocation int       `json:"new_allocation"`
	Successful    *bool     `json:"successful"`
	Archived      string    `gorm:"type:enum('true','false');default:'false'" json:"archived"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (ServerTransfer) TableName() string { return "featherpanel_server_transfers" }

type Subuser struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      int       `gorm:"index" json:"user_id"`
	ServerID    int       `gorm:"index" json:"server_id"`
	Permissions string    `gorm:"type:json" json:"permissions"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Subuser) TableName() string { return "featherpanel_server_subusers" }
