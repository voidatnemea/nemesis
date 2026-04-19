package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mythicalsystems/nemesis-backend/internal/database"
	"github.com/mythicalsystems/nemesis-backend/internal/models"
	"github.com/mythicalsystems/nemesis-backend/pkg/utils"
)

type LocationsController struct{}

func (l *LocationsController) Index(c *gin.Context) {
	var locations []models.Location
	database.DB.Find(&locations)
	utils.Success(c, locations, "Locations retrieved", http.StatusOK)
}

func (l *LocationsController) Show(c *gin.Context) {
	id := c.Param("id")
	var location models.Location
	if err := database.DB.First(&location, id).Error; err != nil {
		utils.Error(c, "Location not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	utils.Success(c, location, "Location retrieved", http.StatusOK)
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
