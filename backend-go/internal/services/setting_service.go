package services

import (
	"github.com/mythicalsystems/nemesis-backend/internal/database"
	"github.com/mythicalsystems/nemesis-backend/internal/models"
)

func GetSetting(key string, fallback string) string {
	var setting models.Setting
	if err := database.DB.Where("`key` = ?", key).First(&setting).Error; err != nil {
		return fallback
	}
	return setting.Value
}

func SetSetting(key, value string) {
	var setting models.Setting
	if err := database.DB.Where("`key` = ?", key).First(&setting).Error; err != nil {
		setting = models.Setting{Key: key, Value: value}
		database.DB.Create(&setting)
	} else {
		database.DB.Model(&setting).Update("value", value)
	}
}
