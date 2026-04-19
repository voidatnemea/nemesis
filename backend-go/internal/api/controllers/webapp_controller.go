package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/mythicalsystems/nemesis-backend/internal/services"
)

type WebAppController struct{}

func (w *WebAppController) Index(c *gin.Context) {
	appName := services.GetSetting("app:name", "FeatherPanel")
	appLogoWhite := services.GetSetting("app:logo_white", "https://github.com/featherpanel-com.png")
	appLogoDark := services.GetSetting("app:logo_dark", "https://github.com/featherpanel-com.png")

	manifest := gin.H{
		"name":        appName,
		"short_name":  appName,
		"description": appName + "'s game server management dashboard",
		"start_url":   "/",
		"scope":       "/",
		"display":     "standalone",
		"theme_color": "#6366f1",
		"background_color": "#020203",
		"orientation":      "any",
		"categories":       []string{"utilities", "productivity", "server", "management", "games"},
		"lang":             "en",
		"dir":              "ltr",
		"id":               "/api/manifest.webmanifest",
		"prefer_related_applications": false,
		"icons": []gin.H{
			{"src": appLogoWhite, "sizes": "72x72", "type": "image/png", "purpose": "any"},
			{"src": appLogoWhite, "sizes": "96x96", "type": "image/png", "purpose": "any"},
			{"src": appLogoWhite, "sizes": "128x128", "type": "image/png", "purpose": "any"},
			{"src": appLogoWhite, "sizes": "144x144", "type": "image/png", "purpose": "any"},
			{"src": appLogoWhite, "sizes": "152x152", "type": "image/png", "purpose": "any"},
			{"src": appLogoWhite, "sizes": "192x192", "type": "image/png", "purpose": "any"},
			{"src": appLogoWhite, "sizes": "384x384", "type": "image/png", "purpose": "any"},
			{"src": appLogoWhite, "sizes": "512x512", "type": "image/png", "purpose": "any"},
			{"src": appLogoDark, "sizes": "192x192", "type": "image/png", "purpose": "any dark"},
			{"src": appLogoDark, "sizes": "512x512", "type": "image/png", "purpose": "any dark"},
			{"src": appLogoWhite, "sizes": "192x192", "type": "image/png", "purpose": "maskable"},
			{"src": appLogoWhite, "sizes": "512x512", "type": "image/png", "purpose": "maskable"},
		},
		"shortcuts": []gin.H{
			{
				"name":       "Dashboard",
				"short_name": "Dashboard",
				"description": "View your game servers",
				"url":        "/dashboard",
				"icons": []gin.H{
					{"src": appLogoWhite, "sizes": "192x192", "type": "image/png"},
				},
			},
			{
				"name":       "Account Settings",
				"short_name": "Account",
				"description": "Manage your account",
				"url":        "/dashboard/account",
				"icons": []gin.H{
					{"src": appLogoWhite, "sizes": "192x192", "type": "image/png"},
				},
			},
		},
		"display_override":     []string{"window-controls-overlay", "standalone", "minimal-ui", "browser"},
		"related_applications": []interface{}{},
	}

	c.JSON(200, manifest)
}
