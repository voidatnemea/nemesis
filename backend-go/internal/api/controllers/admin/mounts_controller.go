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
	p := utils.GetPagination(c)
	var mounts []models.Mount
	var total int64
	database.DB.Model(&models.Mount{}).Count(&total)
	database.DB.Offset(p.Offset).Limit(p.PerPage).Find(&mounts)
	var nodes []models.Node
	database.DB.Find(&nodes)
	var spells []models.Spell
	database.DB.Find(&spells)
	utils.Success(c, gin.H{
		"mounts":     mounts,
		"pagination": utils.BuildPagination(p, total),
		"nodes":      nodes,
		"spells":     spells,
	}, "Mounts retrieved", http.StatusOK)
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
	utils.Success(c, gin.H{"mount_id": mount.ID, "mount": mount}, "Mount created", http.StatusCreated)
}

func (m *MountsController) UpdateLinks(c *gin.Context) {
	id := c.Param("id")
	var mount models.Mount
	if err := database.DB.First(&mount, id).Error; err != nil {
		utils.Error(c, "Mount not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	var req struct {
		Nodes  []int `json:"nodes"`
		Spells []int `json:"spells"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	utils.Success(c, mount, "Mount links updated", http.StatusOK)
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
