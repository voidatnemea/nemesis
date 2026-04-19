package models

import "time"

type TicketCategory struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string    `gorm:"type:varchar(255)" json:"name"`
	Icon         string    `gorm:"type:varchar(255)" json:"icon"`
	Color        string    `gorm:"type:varchar(50)" json:"color"`
	SupportEmail string    `gorm:"type:varchar(255)" json:"support_email"`
	OpenHours    string    `gorm:"type:varchar(255)" json:"open_hours"`
	Description  string    `gorm:"type:text" json:"description"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (TicketCategory) TableName() string { return "featherpanel_ticket_categories" }

type TicketStatus struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(255)" json:"name"`
	Color     string    `gorm:"type:varchar(50)" json:"color"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (TicketStatus) TableName() string { return "featherpanel_ticket_statuses" }

type TicketPriority struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(255)" json:"name"`
	Color     string    `gorm:"type:varchar(50)" json:"color"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (TicketPriority) TableName() string { return "featherpanel_ticket_priorities" }

type Ticket struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID        string    `gorm:"uniqueIndex;type:varchar(36)" json:"uuid"`
	Title       string    `gorm:"type:varchar(255)" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	UserUUID    string    `gorm:"type:varchar(36);index" json:"user_uuid"`
	ServerID    *int      `gorm:"index" json:"server_id"`
	CategoryID  *int      `json:"category_id"`
	StatusID    *int      `json:"status_id"`
	PriorityID  *int      `json:"priority_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Ticket) TableName() string { return "featherpanel_tickets" }

type TicketMessage struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	TicketID  int       `gorm:"index" json:"ticket_id"`
	UserUUID  *string   `gorm:"type:varchar(36)" json:"user_uuid"`
	Message   string    `gorm:"type:text" json:"message"`
	IsStaff   string    `gorm:"type:enum('true','false');default:'false'" json:"is_staff"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (TicketMessage) TableName() string { return "featherpanel_ticket_messages" }

type TicketAttachment struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	TicketID  int       `gorm:"index" json:"ticket_id"`
	MessageID *int      `json:"message_id"`
	FileName  string    `gorm:"type:varchar(255)" json:"file_name"`
	FilePath  string    `gorm:"type:varchar(500)" json:"file_path"`
	FileSize  int64     `json:"file_size"`
	FileType  string    `gorm:"type:varchar(100)" json:"file_type"`
	CreatedAt time.Time `json:"created_at"`
}

func (TicketAttachment) TableName() string { return "featherpanel_ticket_attachments" }
