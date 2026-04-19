package models

import (
	"gorm.io/gorm"
)

type Setting struct {
	gorm.Model
	Key   string `gorm:"uniqueIndex"`
	Value string
}

func (Setting) TableName() string {
	return "featherpanel_settings"
}
