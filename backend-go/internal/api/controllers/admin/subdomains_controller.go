package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mythicalsystems/nemesis-backend/internal/database"
	"github.com/mythicalsystems/nemesis-backend/internal/models"
	"github.com/mythicalsystems/nemesis-backend/pkg/utils"
)

type AdminSubdomainsController struct{}

func (s *AdminSubdomainsController) Index(c *gin.Context) {
	var domains []models.SubdomainDomain
	database.DB.Find(&domains)
	utils.Success(c, domains, "Subdomain domains retrieved", http.StatusOK)
}

func (s *AdminSubdomainsController) Show(c *gin.Context) {
	uuid := c.Param("uuid")
	var domain models.SubdomainDomain
	if err := database.DB.Where("uuid = ?", uuid).First(&domain).Error; err != nil {
		utils.Error(c, "Domain not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	utils.Success(c, domain, "Domain retrieved", http.StatusOK)
}

func (s *AdminSubdomainsController) Create(c *gin.Context) {
	var req struct {
		Domain              string `json:"domain" binding:"required"`
		CloudflareAccountID string `json:"cloudflare_account_id" binding:"required"`
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
	database.DB.Create(&domain)
	utils.Success(c, domain, "Domain created", http.StatusCreated)
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
	utils.Success(c, domain, "Domain updated", http.StatusOK)
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
	utils.Success(c, subdomains, "Subdomains retrieved", http.StatusOK)
}
