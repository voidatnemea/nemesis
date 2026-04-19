package admin

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mythicalsystems/nemesis-backend/internal/database"
	"github.com/mythicalsystems/nemesis-backend/internal/models"
	"github.com/mythicalsystems/nemesis-backend/pkg/utils"
)

type SettingsController struct{}

func (s *SettingsController) Index(c *gin.Context) {
	var settings []models.Setting
	database.DB.Find(&settings)
	utils.Success(c, settings, "Settings retrieved", http.StatusOK)
}

func (s *SettingsController) Categories(c *gin.Context) {
	var settings []models.Setting
	database.DB.Find(&settings)
	categories := map[string]bool{}
	for _, setting := range settings {
		parts := strings.SplitN(setting.Key, ":", 2)
		if len(parts) == 2 {
			categories[parts[0]] = true
		}
	}
	cats := make([]string, 0, len(categories))
	for k := range categories {
		cats = append(cats, k)
	}
	utils.Success(c, cats, "Categories retrieved", http.StatusOK)
}

func (s *SettingsController) GetByCategory(c *gin.Context) {
	category := c.Param("category")
	var settings []models.Setting
	database.DB.Where("`key` LIKE ?", category+":%").Find(&settings)
	utils.Success(c, settings, "Settings retrieved", http.StatusOK)
}

func (s *SettingsController) Show(c *gin.Context) {
	key := c.Param("setting")
	var setting models.Setting
	if err := database.DB.Where("`key` = ?", key).First(&setting).Error; err != nil {
		utils.Error(c, "Setting not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	utils.Success(c, setting, "Setting retrieved", http.StatusOK)
}

func (s *SettingsController) Update(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	for key, value := range req {
		var setting models.Setting
		if database.DB.Where("`key` = ?", key).First(&setting).Error != nil {
			database.DB.Create(&models.Setting{Key: key, Value: value})
		} else {
			database.DB.Model(&setting).Update("value", value)
		}
	}
	utils.Success(c, nil, "Settings updated", http.StatusOK)
}
