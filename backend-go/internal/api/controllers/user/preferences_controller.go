package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mythicalsystems/nemesis-backend/internal/database"
	"github.com/mythicalsystems/nemesis-backend/internal/models"
	"github.com/mythicalsystems/nemesis-backend/pkg/utils"
)

type UserPreferencesController struct{}

func (p *UserPreferencesController) Get(c *gin.Context) {
	ctxUser, _ := c.Get("user")
	user := ctxUser.(models.User)

	var prefs []models.UserPreference
	database.DB.Where("user_uuid = ?", user.UUID).Find(&prefs)

	result := map[string]string{}
	for _, pref := range prefs {
		result[pref.Key] = pref.Value
	}

	utils.Success(c, result, "Preferences retrieved", http.StatusOK)
}

func (p *UserPreferencesController) Update(c *gin.Context) {
	ctxUser, _ := c.Get("user")
	user := ctxUser.(models.User)

	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}

	for key, value := range req {
		var pref models.UserPreference
		result := database.DB.Where("user_uuid = ? AND key = ?", user.UUID, key).First(&pref)
		if result.Error != nil {
			pref = models.UserPreference{UserUUID: user.UUID, Key: key, Value: value}
			database.DB.Create(&pref)
		} else {
			database.DB.Model(&pref).Update("value", value)
		}
	}

	utils.Success(c, req, "Preferences updated", http.StatusOK)
}
