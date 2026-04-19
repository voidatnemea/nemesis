package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mythicalsystems/nemesis-backend/internal/database"
	"github.com/mythicalsystems/nemesis-backend/internal/models"
	"github.com/mythicalsystems/nemesis-backend/pkg/utils"
)

type AllocationsController struct{}

func (a *AllocationsController) Index(c *gin.Context) {
	p := utils.GetPagination(c)
	nodeID := c.Query("node_id")
	var allocs []models.Allocation
	var total int64
	q := database.DB.Model(&models.Allocation{})
	if nodeID != "" {
		q = q.Where("node_id = ?", nodeID)
	}
	q.Count(&total)
	q.Offset(p.Offset).Limit(p.PerPage).Find(&allocs)
	utils.Success(c, gin.H{"data": allocs, "total": total}, "Allocations retrieved", http.StatusOK)
}

func (a *AllocationsController) Show(c *gin.Context) {
	id := c.Param("id")
	var alloc models.Allocation
	if err := database.DB.First(&alloc, id).Error; err != nil {
		utils.Error(c, "Allocation not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	utils.Success(c, alloc, "Allocation retrieved", http.StatusOK)
}

func (a *AllocationsController) Create(c *gin.Context) {
	var req struct {
		NodeID int    `json:"node_id" binding:"required"`
		IP     string `json:"ip" binding:"required"`
		Ports  []int  `json:"ports" binding:"required"`
		Alias  string `json:"alias"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	var created []models.Allocation
	for _, port := range req.Ports {
		alloc := models.Allocation{NodeID: req.NodeID, IP: req.IP, Port: port, Alias: req.Alias}
		database.DB.Create(&alloc)
		created = append(created, alloc)
	}
	utils.Success(c, created, "Allocations created", http.StatusCreated)
}

func (a *AllocationsController) Update(c *gin.Context) {
	id := c.Param("id")
	var alloc models.Allocation
	if err := database.DB.First(&alloc, id).Error; err != nil {
		utils.Error(c, "Allocation not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	var req map[string]interface{}
	c.ShouldBindJSON(&req)
	delete(req, "id")
	database.DB.Model(&alloc).Updates(req)
	utils.Success(c, alloc, "Allocation updated", http.StatusOK)
}

func (a *AllocationsController) Delete(c *gin.Context) {
	id := c.Param("id")
	database.DB.Delete(&models.Allocation{}, id)
	utils.Success(c, nil, "Allocation deleted", http.StatusOK)
}

func (a *AllocationsController) GetAvailable(c *gin.Context) {
	nodeID := c.Query("node_id")
	var allocs []models.Allocation
	q := database.DB.Where("server_id IS NULL")
	if nodeID != "" {
		q = q.Where("node_id = ?", nodeID)
	}
	q.Find(&allocs)
	utils.Success(c, allocs, "Available allocations retrieved", http.StatusOK)
}

func (a *AllocationsController) BulkDelete(c *gin.Context) {
	var req struct {
		IDs []int `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	database.DB.Where("id IN ? AND server_id IS NULL", req.IDs).Delete(&models.Allocation{})
	utils.Success(c, nil, "Allocations deleted", http.StatusOK)
}

func (a *AllocationsController) DeleteUnused(c *gin.Context) {
	nodeID := c.Query("node_id")
	q := database.DB.Where("server_id IS NULL")
	if nodeID != "" {
		q = q.Where("node_id = ?", nodeID)
	}
	q.Delete(&models.Allocation{})
	utils.Success(c, nil, "Unused allocations deleted", http.StatusOK)
}
