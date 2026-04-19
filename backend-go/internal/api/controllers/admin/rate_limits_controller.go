package admin

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mythicalsystems/nemesis-backend/internal/services"
	"github.com/mythicalsystems/nemesis-backend/pkg/utils"
)

type RateLimitsController struct{}

type routeLimit struct {
	Route     string  `json:"route"`
	PerSecond float64 `json:"per_second"`
	PerMinute float64 `json:"per_minute"`
	PerHour   float64 `json:"per_hour"`
	PerDay    float64 `json:"per_day"`
}

var defaultRoutes = []routeLimit{
	{Route: "/api/user/auth/login", PerSecond: 1, PerMinute: 10, PerHour: 60, PerDay: 200},
	{Route: "/api/user/auth/register", PerSecond: 0.5, PerMinute: 5, PerHour: 20, PerDay: 50},
	{Route: "/api/user/auth/forgot-password", PerSecond: 0.2, PerMinute: 3, PerHour: 10, PerDay: 30},
	{Route: "/api/admin/*", PerSecond: 10, PerMinute: 300, PerHour: 3000, PerDay: 30000},
	{Route: "/api/user/*", PerSecond: 5, PerMinute: 150, PerHour: 1500, PerDay: 15000},
}

func getRateLimitConfig() (bool, []routeLimit) {
	enabled := services.GetSetting("rate_limits:global_enabled", "true") == "true"
	raw := services.GetSetting("rate_limits:routes", "")
	var routes []routeLimit
	if raw != "" {
		if err := json.Unmarshal([]byte(raw), &routes); err != nil {
			routes = defaultRoutes
		}
	} else {
		routes = defaultRoutes
	}
	return enabled, routes
}

func (r *RateLimitsController) GetAll(c *gin.Context) {
	enabled, routes := getRateLimitConfig()
	utils.Success(c, gin.H{
		"global_enabled": enabled,
		"routes":         routes,
		"namespaces":     []string{"api", "auth", "admin"},
	}, "Rate limits retrieved", http.StatusOK)
}

func (r *RateLimitsController) SetGlobal(c *gin.Context) {
	var req struct {
		Enabled bool `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	v := "false"
	if req.Enabled {
		v = "true"
	}
	services.SetSetting("rate_limits:global_enabled", v)
	utils.Success(c, gin.H{"global_enabled": req.Enabled}, "Global rate limit updated", http.StatusOK)
}

func (r *RateLimitsController) UpdateRoute(c *gin.Context) {
	var req routeLimit
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	_, routes := getRateLimitConfig()
	found := false
	for i, rt := range routes {
		if rt.Route == req.Route {
			routes[i] = req
			found = true
			break
		}
	}
	if !found {
		routes = append(routes, req)
	}
	b, _ := json.Marshal(routes)
	services.SetSetting("rate_limits:routes", string(b))
	utils.Success(c, req, "Route updated", http.StatusOK)
}

func (r *RateLimitsController) BulkUpdate(c *gin.Context) {
	var req struct {
		Routes []routeLimit `json:"routes" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	b, _ := json.Marshal(req.Routes)
	services.SetSetting("rate_limits:routes", string(b))
	utils.Success(c, gin.H{"routes": req.Routes}, "Routes updated", http.StatusOK)
}
