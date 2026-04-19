package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mythicalsystems/nemesis-backend/internal/database"
	"github.com/mythicalsystems/nemesis-backend/internal/models"
	"github.com/mythicalsystems/nemesis-backend/pkg/utils"
)

type RealmsController struct{}

func (r *RealmsController) Index(c *gin.Context) {
	p := utils.GetPagination(c)
	var realms []models.Realm
	var total int64
	database.DB.Model(&models.Realm{}).Count(&total)
	database.DB.Offset(p.Offset).Limit(p.PerPage).Find(&realms)
	utils.Success(c, gin.H{"realms": realms, "pagination": utils.BuildPagination(p, total)}, "Realms retrieved", http.StatusOK)
}

func (r *RealmsController) Show(c *gin.Context) {
	id := c.Param("id")
	var realm models.Realm
	if err := database.DB.First(&realm, id).Error; err != nil {
		utils.Error(c, "Realm not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	utils.Success(c, realm, "Realm retrieved", http.StatusOK)
}

func (r *RealmsController) Create(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	realm := models.Realm{Name: req.Name, Description: req.Description}
	database.DB.Create(&realm)
	utils.Success(c, realm, "Realm created", http.StatusCreated)
}

func (r *RealmsController) Update(c *gin.Context) {
	id := c.Param("id")
	var realm models.Realm
	if err := database.DB.First(&realm, id).Error; err != nil {
		utils.Error(c, "Realm not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	var req map[string]interface{}
	c.ShouldBindJSON(&req)
	delete(req, "id")
	database.DB.Model(&realm).Updates(req)
	utils.Success(c, realm, "Realm updated", http.StatusOK)
}

func (r *RealmsController) Delete(c *gin.Context) {
	id := c.Param("id")
	database.DB.Delete(&models.Realm{}, id)
	utils.Success(c, nil, "Realm deleted", http.StatusOK)
}
