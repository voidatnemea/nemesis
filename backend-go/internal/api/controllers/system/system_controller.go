package system

import (
	"context"
	"net/http"
	"os"
	"time"

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
		"branding:custom_css", "branding:custom_js",
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

func (s *SystemController) SelfTest(c *gin.Context) {
	mysqlStatus := false
	mysqlMsg := "Failed"
	if database.DB != nil {
		if sqlDB, err := database.DB.DB(); err == nil {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			if err := sqlDB.PingContext(ctx); err == nil {
				mysqlStatus = true
				mysqlMsg = "Successful"
			} else {
				mysqlMsg = err.Error()
			}
		}
	}

	redisStatus := false
	redisMsg := "Failed"
	if database.Redis != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if _, err := database.Redis.Ping(ctx).Result(); err == nil {
			redisStatus = true
			redisMsg = "Successful"
		} else {
			redisMsg = err.Error()
		}
	}

	permissions := map[string]bool{}
	for _, dir := range []string{"/tmp", "/var/tmp"} {
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			permissions[dir] = true
		} else {
			permissions[dir] = false
		}
	}

	status := "ready"
	if !mysqlStatus {
		status = "degraded"
	}

	utils.Success(c, gin.H{
		"status": status,
		"cached": false,
		"checks": gin.H{
			"mysql": gin.H{
				"status":  mysqlStatus,
				"message": mysqlMsg,
			},
			"redis": gin.H{
				"status":  redisStatus,
				"message": redisMsg,
			},
			"permissions": permissions,
		},
	}, "System is ready", http.StatusOK)
}

func (s *SystemController) GetPluginCSS(c *gin.Context) {
	c.Data(http.StatusOK, "text/css; charset=utf-8", []byte(""))
}

func (s *SystemController) GetPluginJS(c *gin.Context) {
	c.Data(http.StatusOK, "application/javascript; charset=utf-8", []byte(""))
}

func (s *SystemController) GetPluginWidgets(c *gin.Context) {
	utils.Success(c, []interface{}{}, "Plugin widgets retrieved", http.StatusOK)
}

func (s *SystemController) GetPluginSidebar(c *gin.Context) {
	utils.Success(c, []interface{}{}, "Plugin sidebar retrieved", http.StatusOK)
}

func (s *SystemController) GetTranslationLanguages(c *gin.Context) {
	utils.Success(c, []string{"en"}, "Languages retrieved", http.StatusOK)
}

func (s *SystemController) GetTranslation(c *gin.Context) {
	utils.Success(c, gin.H{}, "Translations retrieved", http.StatusOK)
}

type NotificationsController struct{}

func (n *NotificationsController) GetUserNotifications(c *gin.Context) {
	var notifs []models.Notification
	database.DB.Order("created_at DESC").Limit(20).Find(&notifs)
	utils.Success(c, notifs, "Notifications retrieved", http.StatusOK)
}
