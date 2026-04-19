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
