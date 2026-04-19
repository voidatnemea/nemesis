package system

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mythicalsystems/nemesis-backend/internal/database"
	"github.com/mythicalsystems/nemesis-backend/internal/models"
	"github.com/mythicalsystems/nemesis-backend/internal/services"
	"github.com/mythicalsystems/nemesis-backend/pkg/utils"
)

type SystemController struct{}

func (s *SystemController) GetSettings(c *gin.Context) {
	publicKeys := []string{
		"app:name", "app:logo_white", "app:logo_dark",
		"app:description", "app:url", "auth:registration",
		"auth:discord_enabled", "auth:oidc_enabled",
	}
	result := map[string]string{}
	for _, key := range publicKeys {
		result[key] = services.GetSetting(key, "")
	}
	utils.Success(c, result, "Settings retrieved", http.StatusOK)
}

func (s *SystemController) GetOidcProviders(c *gin.Context) {
	var providers []models.OidcProvider
	database.DB.Where("enabled = ?", "true").Find(&providers)
	// Strip client secret from response
	type PublicProvider struct {
		UUID      string `json:"uuid"`
		Name      string `json:"name"`
		IssuerURL string `json:"issuer_url"`
		ClientID  string `json:"client_id"`
	}
	var public []PublicProvider
	for _, p := range providers {
		public = append(public, PublicProvider{
			UUID:      p.UUID,
			Name:      p.Name,
			IssuerURL: p.IssuerURL,
			ClientID:  p.ClientID,
		})
	}
	utils.Success(c, public, "OIDC providers retrieved", http.StatusOK)
}

type NotificationsController struct{}

func (n *NotificationsController) GetUserNotifications(c *gin.Context) {
	var notifs []models.Notification
	database.DB.Order("created_at DESC").Limit(20).Find(&notifs)
	utils.Success(c, notifs, "Notifications retrieved", http.StatusOK)
}
