package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mythicalsystems/nemesis-backend/internal/database"
	"github.com/mythicalsystems/nemesis-backend/internal/models"
	"github.com/mythicalsystems/nemesis-backend/pkg/utils"
)

type MountsController struct{}

func (m *MountsController) Index(c *gin.Context) {
	var mounts []models.Mount
	database.DB.Find(&mounts)
	utils.Success(c, mounts, "Mounts retrieved", http.StatusOK)
}

func (m *MountsController) Show(c *gin.Context) {
	id := c.Param("id")
	var mount models.Mount
	if err := database.DB.First(&mount, id).Error; err != nil {
		utils.Error(c, "Mount not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	utils.Success(c, mount, "Mount retrieved", http.StatusOK)
}

func (m *MountsController) Create(c *gin.Context) {
	var req struct {
		Name          string `json:"name" binding:"required"`
		Description   string `json:"description"`
		Source        string `json:"source" binding:"required"`
		Target        string `json:"target" binding:"required"`
		ReadOnly      string `json:"read_only"`
		UserMountable string `json:"user_mountable"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	mount := models.Mount{
		Name: req.Name, Description: req.Description, Source: req.Source,
		Target: req.Target, ReadOnly: req.ReadOnly, UserMountable: req.UserMountable,
	}
	database.DB.Create(&mount)
	utils.Success(c, mount, "Mount created", http.StatusCreated)
}

func (m *MountsController) Update(c *gin.Context) {
	id := c.Param("id")
	var mount models.Mount
	if err := database.DB.First(&mount, id).Error; err != nil {
		utils.Error(c, "Mount not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	var req map[string]interface{}
	c.ShouldBindJSON(&req)
	delete(req, "id")
	database.DB.Model(&mount).Updates(req)
	utils.Success(c, mount, "Mount updated", http.StatusOK)
}

func (m *MountsController) Delete(c *gin.Context) {
	id := c.Param("id")
	database.DB.Delete(&models.Mount{}, id)
	utils.Success(c, nil, "Mount deleted", http.StatusOK)
}
