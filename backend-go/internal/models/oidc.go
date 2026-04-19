package models

import "time"

type OidcProvider struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID         string    `gorm:"uniqueIndex;type:varchar(36)" json:"uuid"`
	Name         string    `gorm:"type:varchar(255)" json:"name"`
	IssuerURL    string    `gorm:"type:varchar(500)" json:"issuer_url"`
	ClientID     string    `gorm:"type:varchar(255)" json:"client_id"`
	ClientSecret string    `gorm:"type:varchar(500)" json:"-"`
	RedirectURI  string    `gorm:"type:varchar(500)" json:"redirect_uri"`
	Scopes       string    `gorm:"type:varchar(500);default:'openid email profile'" json:"scopes"`
	Enabled      string    `gorm:"type:enum('true','false');default:'true'" json:"enabled"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (OidcProvider) TableName() string { return "featherpanel_oidc_providers" }

type SsoToken struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserUUID  string    `gorm:"type:varchar(36);index" json:"user_uuid"`
	Token     string    `gorm:"uniqueIndex;type:varchar(255)" json:"token"`
	Used      string    `gorm:"type:enum('true','false');default:'false'" json:"used"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

func (SsoToken) TableName() string { return "featherpanel_sso_tokens" }
