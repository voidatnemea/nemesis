package admin

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mythicalsystems/nemesis-backend/internal/database"
	"github.com/mythicalsystems/nemesis-backend/internal/models"
	"github.com/mythicalsystems/nemesis-backend/internal/services"
	"github.com/mythicalsystems/nemesis-backend/pkg/utils"
)

type AdminSubdomainsController struct{}

func (s *AdminSubdomainsController) Index(c *gin.Context) {
	p := utils.GetPagination(c)
	search := strings.TrimSpace(c.Query("search"))

	q := database.DB.Model(&models.SubdomainDomain{})
	if search != "" {
		q = q.Where("domain LIKE ?", "%"+search+"%")
	}
	var total int64
	q.Count(&total)

	var domains []models.SubdomainDomain
	q.Order("created_at DESC").Offset(p.Offset).Limit(p.PerPage).Find(&domains)
	if domains == nil {
		domains = []models.SubdomainDomain{}
	}

	utils.Success(c, gin.H{
		"domains":    domains,
		"pagination": utils.BuildPagination(p, total),
	}, "Subdomain domains retrieved", http.StatusOK)
}

func (s *AdminSubdomainsController) Show(c *gin.Context) {
	uuid := c.Param("uuid")
	var domain models.SubdomainDomain
	if err := database.DB.Where("uuid = ?", uuid).First(&domain).Error; err != nil {
		utils.Error(c, "Domain not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	utils.Success(c, gin.H{"domain": domain}, "Domain retrieved", http.StatusOK)
}

func (s *AdminSubdomainsController) Create(c *gin.Context) {
	var req struct {
		Domain              string `json:"domain" binding:"required"`
		CloudflareAccountID string `json:"cloudflare_account_id"`
		CloudflareZoneID    string `json:"cloudflare_zone_id"`
		CloudflareAPIToken  string `json:"cloudflare_api_token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	domain := models.SubdomainDomain{
		UUID:                utils.GenerateUUID(),
		Domain:              req.Domain,
		CloudflareAccountID: req.CloudflareAccountID,
		CloudflareZoneID:    req.CloudflareZoneID,
		CloudflareAPIToken:  req.CloudflareAPIToken,
	}
	if err := database.DB.Create(&domain).Error; err != nil {
		utils.Error(c, "Failed to create domain", "DATABASE_ERROR", http.StatusConflict, nil)
		return
	}
	utils.Success(c, gin.H{"domain": domain}, "Domain created", http.StatusCreated)
}

func (s *AdminSubdomainsController) Update(c *gin.Context) {
	uuid := c.Param("uuid")
	var domain models.SubdomainDomain
	if err := database.DB.Where("uuid = ?", uuid).First(&domain).Error; err != nil {
		utils.Error(c, "Domain not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	var req map[string]interface{}
	c.ShouldBindJSON(&req)
	delete(req, "id")
	delete(req, "uuid")
	database.DB.Model(&domain).Updates(req)
	utils.Success(c, gin.H{"domain": domain}, "Domain updated", http.StatusOK)
}

func (s *AdminSubdomainsController) Delete(c *gin.Context) {
	uuid := c.Param("uuid")
	database.DB.Where("uuid = ?", uuid).Delete(&models.SubdomainDomain{})
	utils.Success(c, nil, "Domain deleted", http.StatusOK)
}

func (s *AdminSubdomainsController) ListSubdomains(c *gin.Context) {
	uuid := c.Param("uuid")
	var domain models.SubdomainDomain
	if err := database.DB.Where("uuid = ?", uuid).First(&domain).Error; err != nil {
		utils.Error(c, "Domain not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	var subdomains []models.Subdomain
	database.DB.Where("domain_id = ?", domain.ID).Find(&subdomains)
	if subdomains == nil {
		subdomains = []models.Subdomain{}
	}
	utils.Success(c, gin.H{"subdomains": subdomains}, "Subdomains retrieved", http.StatusOK)
}

func (s *AdminSubdomainsController) GetSettings(c *gin.Context) {
	email := services.GetSetting("subdomains:cloudflare_email", "")
	apiKeySet := services.GetSetting("subdomains:cloudflare_api_key", "") != ""
	maxPerServer, _ := strconv.Atoi(services.GetSetting("subdomains:max_subdomains_per_server", "1"))
	if maxPerServer <= 0 {
		maxPerServer = 1
	}
	allow := services.GetSetting("subdomains:allow_user_subdomains", "false") == "true"

	utils.Success(c, gin.H{
		"settings": gin.H{
			"cloudflare_email":          email,
			"cloudflare_api_key_set":    apiKeySet,
			"max_subdomains_per_server": maxPerServer,
			"allow_user_subdomains":     allow,
		},
	}, "Subdomain settings retrieved", http.StatusOK)
}

func (s *AdminSubdomainsController) UpdateSettings(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	for k, v := range req {
		switch k {
		case "cloudflare_email":
			if sv, ok := v.(string); ok {
				services.SetSetting("subdomains:cloudflare_email", sv)
			}
		case "cloudflare_api_key":
			if sv, ok := v.(string); ok && sv != "" {
				services.SetSetting("subdomains:cloudflare_api_key", sv)
			}
		case "max_subdomains_per_server":
			switch x := v.(type) {
			case float64:
				services.SetSetting("subdomains:max_subdomains_per_server", formatNum(int(x)))
			case string:
				services.SetSetting("subdomains:max_subdomains_per_server", x)
			}
		case "allow_user_subdomains":
			if b, ok := v.(bool); ok {
				if b {
					services.SetSetting("subdomains:allow_user_subdomains", "true")
				} else {
					services.SetSetting("subdomains:allow_user_subdomains", "false")
				}
			}
		}
	}
	utils.Success(c, nil, "Subdomain settings updated", http.StatusOK)
}

func (s *AdminSubdomainsController) ListSpells(c *gin.Context) {
	type spellRow struct {
		UUID string `json:"uuid"`
		Name string `json:"name"`
	}
	var rows []spellRow
	database.DB.Raw("SELECT uuid, name FROM featherpanel_spells WHERE deleted = 'false' ORDER BY name ASC").Scan(&rows)
	if rows == nil {
		rows = []spellRow{}
	}
	utils.Success(c, gin.H{"spells": rows}, "Subdomain spells retrieved", http.StatusOK)
}

func formatNum(n int) string {
	if n <= 0 {
		return "0"
	}
	return strconv.Itoa(n)
}
