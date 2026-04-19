package middleware

import (
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mythicalsystems/nemesis-backend/internal/database"
	"github.com/mythicalsystems/nemesis-backend/internal/models"
	"github.com/mythicalsystems/nemesis-backend/pkg/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		var authType string

		// 1. Try Session (Cookie)
		cookie, err := (c).Cookie("remember_token")
		if err == nil && cookie != "" {
			if err := database.DB.Where("remember_token = ?", utils.HashToken(cookie)).First(&user).Error; err == nil {
				if user.Banned == "true" {
					utils.Error(c, "User is banned", "USER_BANNED", http.StatusForbidden, nil)
					c.Abort()
					return
				}
				authType = "session"
			}
		}

		// 2. Try API Key (Bearer Token)
		if authType == "" {
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				token := strings.TrimPrefix(authHeader, "Bearer ")

				// Try public key first, then private key (matches PHP backend behaviour)
				var apiClient models.ApiClient
				if database.DB.Where("public_key = ?", token).First(&apiClient).Error != nil {
					database.DB.Where("private_key = ?", token).First(&apiClient)
				}
				if apiClient.ID != 0 {
					if !isIPAllowed(apiClient.AllowedIps, c.ClientIP()) {
						utils.Error(c, "This API key cannot be used from your IP address", "API_KEY_IP_NOT_ALLOWED", http.StatusForbidden, nil)
						c.Abort()
						return
					}
					if err := database.DB.Where("uuid = ?", apiClient.UserUUID).First(&user).Error; err == nil {
						if user.Banned == "true" {
							utils.Error(c, "User is banned", "USER_BANNED", http.StatusForbidden, nil)
							c.Abort()
							return
						}
						authType = "api_key"
						c.Set("api_client", apiClient)
					}
				} else {
					utils.Error(c, "Invalid API key", "INVALID_API_KEY", http.StatusUnauthorized, nil)
					c.Abort()
					return
				}
			}
		}

		if authType == "" {
			utils.Error(c, "You are not allowed to access this resource!", "INVALID_ACCOUNT_TOKEN", http.StatusUnauthorized, nil)
			c.Abort()
			return
		}

		// Set context variables
		c.Set("user", user)
		c.Set("auth_type", authType)

		c.Next()
	}
}

// isIPAllowed checks whether the request IP is permitted when the API client's
// AllowedIps column is populated. AllowedIps is expected to be a comma-separated
// list of IPs or CIDR blocks.
func isIPAllowed(allowed string, requestIP string) bool {
	if strings.TrimSpace(allowed) == "" {
		return true
	}
	ip := net.ParseIP(requestIP)
	if ip == nil {
		return false
	}
	for _, entry := range strings.Split(allowed, ",") {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		if parsed := net.ParseIP(entry); parsed != nil {
			if parsed.Equal(ip) {
				return true
			}
			continue
		}
		if _, cidr, err := net.ParseCIDR(entry); err == nil {
			if cidr.Contains(ip) {
				return true
			}
		}
	}
	return false
}
