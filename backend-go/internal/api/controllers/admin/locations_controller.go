package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mythicalsystems/nemesis-backend/internal/database"
	"github.com/mythicalsystems/nemesis-backend/internal/models"
	"github.com/mythicalsystems/nemesis-backend/pkg/utils"
)

type LocationsController struct{}

var countryCodes = map[string]string{
	"AF": "Afghanistan", "AL": "Albania", "DZ": "Algeria", "AR": "Argentina", "AU": "Australia",
	"AT": "Austria", "BE": "Belgium", "BR": "Brazil", "CA": "Canada", "CL": "Chile",
	"CN": "China", "CO": "Colombia", "HR": "Croatia", "CZ": "Czech Republic", "DK": "Denmark",
	"EG": "Egypt", "FI": "Finland", "FR": "France", "DE": "Germany", "GR": "Greece",
	"HK": "Hong Kong", "HU": "Hungary", "IN": "India", "ID": "Indonesia", "IE": "Ireland",
	"IL": "Israel", "IT": "Italy", "JP": "Japan", "KR": "South Korea", "MX": "Mexico",
	"NL": "Netherlands", "NZ": "New Zealand", "NO": "Norway", "PL": "Poland", "PT": "Portugal",
	"RO": "Romania", "RU": "Russia", "SA": "Saudi Arabia", "SG": "Singapore", "ZA": "South Africa",
	"ES": "Spain", "SE": "Sweden", "CH": "Switzerland", "TW": "Taiwan", "TR": "Turkey",
	"UA": "Ukraine", "GB": "United Kingdom", "US": "United States", "VN": "Vietnam",
}

func (l *LocationsController) Index(c *gin.Context) {
	p := utils.GetPagination(c)
	var locations []models.Location
	var total int64
	database.DB.Model(&models.Location{}).Count(&total)
	database.DB.Offset(p.Offset).Limit(p.PerPage).Find(&locations)
	utils.Success(c, gin.H{
		"locations":    locations,
		"pagination":   utils.BuildPagination(p, total),
		"country_codes": countryCodes,
	}, "Locations retrieved", http.StatusOK)
}

func (l *LocationsController) Show(c *gin.Context) {
	id := c.Param("id")
	var location models.Location
	if err := database.DB.First(&location, id).Error; err != nil {
		utils.Error(c, "Location not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	utils.Success(c, gin.H{"location": location}, "Location retrieved", http.StatusOK)
}

func (l *LocationsController) Create(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		FlagCode    string `json:"flag_code"`
		Type        string `json:"type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	if req.Type == "" {
		req.Type = "game"
	}
	loc := models.Location{Name: req.Name, Description: req.Description, FlagCode: req.FlagCode, Type: req.Type}
	database.DB.Create(&loc)
	utils.Success(c, loc, "Location created", http.StatusCreated)
}

func (l *LocationsController) Update(c *gin.Context) {
	id := c.Param("id")
	var loc models.Location
	if err := database.DB.First(&loc, id).Error; err != nil {
		utils.Error(c, "Location not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	var req map[string]interface{}
	c.ShouldBindJSON(&req)
	delete(req, "id")
	database.DB.Model(&loc).Updates(req)
	utils.Success(c, loc, "Location updated", http.StatusOK)
}

func (l *LocationsController) Delete(c *gin.Context) {
	id := c.Param("id")
	database.DB.Delete(&models.Location{}, id)
	utils.Success(c, nil, "Location deleted", http.StatusOK)
}
