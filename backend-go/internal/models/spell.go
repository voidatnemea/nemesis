package models

import "time"

type Spell struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UUID        string    `gorm:"uniqueIndex;type:varchar(36)" json:"uuid"`
	RealmID     int       `gorm:"index" json:"realm_id"`
	Author      string    `gorm:"type:varchar(255)" json:"author"`
	Name        string    `gorm:"type:varchar(255)" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	DockerImages string   `gorm:"type:text" json:"docker_images"`
	Startup     string    `gorm:"type:text" json:"startup"`
	Config      string    `gorm:"type:longtext" json:"config"`
	Scripts     string    `gorm:"type:longtext" json:"scripts"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Spell) TableName() string { return "featherpanel_spells" }

type SpellVariable struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	SpellID      int       `gorm:"index" json:"spell_id"`
	Name         string    `gorm:"type:varchar(255)" json:"name"`
	Description  string    `gorm:"type:text" json:"description"`
	EnvVariable  string    `gorm:"type:varchar(255)" json:"env_variable"`
	DefaultValue string    `gorm:"type:text" json:"default_value"`
	UserViewable string    `gorm:"type:enum('true','false');default:'false'" json:"user_viewable"`
	UserEditable string    `gorm:"type:enum('true','false');default:'false'" json:"user_editable"`
	Rules        string    `gorm:"type:text" json:"rules"`
	FieldType    string    `gorm:"type:varchar(50);default:'text'" json:"field_type"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (SpellVariable) TableName() string { return "featherpanel_spell_variables" }
