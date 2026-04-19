package models

import "time"

type SubdomainDomain struct {
	ID                   uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID                 string    `gorm:"uniqueIndex;type:varchar(36)" json:"uuid"`
	Domain               string    `gorm:"type:varchar(255)" json:"domain"`
	CloudflareAccountID  string    `gorm:"type:varchar(255)" json:"cloudflare_account_id"`
	CloudflareZoneID     string    `gorm:"type:varchar(255)" json:"cloudflare_zone_id"`
	CloudflareAPIToken   string    `gorm:"type:varchar(255)" json:"-"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

func (SubdomainDomain) TableName() string { return "featherpanel_subdomain_manager_domains" }

type Subdomain struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID      string    `gorm:"uniqueIndex;type:varchar(36)" json:"uuid"`
	ServerID  int       `gorm:"index" json:"server_id"`
	DomainID  int       `gorm:"index" json:"domain_id"`
	Subdomain string    `gorm:"type:varchar(255)" json:"subdomain"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Subdomain) TableName() string { return "featherpanel_subdomain_manager_subdomains" }
