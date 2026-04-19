package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mythicalsystems/nemesis-backend/internal/database"
	"github.com/mythicalsystems/nemesis-backend/internal/models"
	"github.com/mythicalsystems/nemesis-backend/pkg/utils"
)

type ServersController struct{}

func (s *ServersController) Index(c *gin.Context) {
	p := utils.GetPagination(c)
	search := c.Query("search")
	var servers []models.Server
	var total int64
	q := database.DB.Model(&models.Server{})
	if search != "" {
		q = q.Where("name LIKE ? OR uuid LIKE ? OR uuid_short LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	q.Count(&total)
	q.Offset(p.Offset).Limit(p.PerPage).Find(&servers)
	utils.Success(c, gin.H{"data": servers, "total": total, "page": p.Page, "per_page": p.PerPage}, "Servers retrieved", http.StatusOK)
}

func (s *ServersController) Show(c *gin.Context) {
	id := c.Param("id")
	var server models.Server
	if err := database.DB.Where("id = ? OR uuid = ?", id, id).First(&server).Error; err != nil {
		utils.Error(c, "Server not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	utils.Success(c, server, "Server retrieved", http.StatusOK)
}

func (s *ServersController) ShowByExternalID(c *gin.Context) {
	externalID := c.Param("externalId")
	var server models.Server
	if err := database.DB.Where("external_id = ?", externalID).First(&server).Error; err != nil {
		utils.Error(c, "Server not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	utils.Success(c, server, "Server retrieved", http.StatusOK)
}

func (s *ServersController) Create(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		OwnerID     int    `json:"owner_id" binding:"required"`
		NodeID      int    `json:"node_id" binding:"required"`
		Memory      int64  `json:"memory" binding:"required"`
		Swap        int64  `json:"swap"`
		Disk        int64  `json:"disk" binding:"required"`
		IO          int    `json:"io"`
		CPU         int    `json:"cpu" binding:"required"`
		AllocationID int   `json:"allocation_id" binding:"required"`
		RealmsID    int    `json:"realms_id" binding:"required"`
		SpellID     int    `json:"spell_id" binding:"required"`
		Image       string `json:"image" binding:"required"`
		Startup     string `json:"startup" binding:"required"`
		DatabaseLimit   int `json:"database_limit"`
		AllocationLimit int `json:"allocation_limit"`
		BackupLimit     int `json:"backup_limit"`
		ExternalID  string `json:"external_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	if req.IO == 0 {
		req.IO = 500
	}
	serverUUID := utils.GenerateUUID()
	uuidShort := serverUUID[:8]
	server := models.Server{
		UUID:            serverUUID,
		UUIDShort:       uuidShort,
		Name:            req.Name,
		Description:     req.Description,
		OwnerID:         req.OwnerID,
		NodeID:          req.NodeID,
		Memory:          req.Memory,
		Swap:            req.Swap,
		Disk:            req.Disk,
		IO:              req.IO,
		CPU:             req.CPU,
		AllocationID:    req.AllocationID,
		RealmsID:        req.RealmsID,
		SpellID:         req.SpellID,
		Image:           req.Image,
		Startup:         req.Startup,
		Status:          "installing",
		DatabaseLimit:   req.DatabaseLimit,
		AllocationLimit: req.AllocationLimit,
		BackupLimit:     req.BackupLimit,
		ExternalID:      req.ExternalID,
	}
	if err := database.DB.Create(&server).Error; err != nil {
		utils.Error(c, "Failed to create server", "DATABASE_ERROR", http.StatusInternalServerError, nil)
		return
	}
	// Mark allocation as used
	database.DB.Model(&models.Allocation{}).Where("id = ?", req.AllocationID).Update("server_id", server.ID)
	utils.Success(c, server, "Server created", http.StatusCreated)
}

func (s *ServersController) Update(c *gin.Context) {
	id := c.Param("id")
	var server models.Server
	if err := database.DB.Where("id = ? OR uuid = ?", id, id).First(&server).Error; err != nil {
		utils.Error(c, "Server not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	delete(req, "id")
	delete(req, "uuid")
	delete(req, "uuid_short")
	database.DB.Model(&server).Updates(req)
	utils.Success(c, server, "Server updated", http.StatusOK)
}

func (s *ServersController) Delete(c *gin.Context) {
	id := c.Param("id")
	database.DB.Where("id = ? OR uuid = ?", id, id).Delete(&models.Server{})
	utils.Success(c, nil, "Server deleted", http.StatusOK)
}

func (s *ServersController) HardDelete(c *gin.Context) {
	id := c.Param("id")
	database.DB.Unscoped().Where("id = ? OR uuid = ?", id, id).Delete(&models.Server{})
	utils.Success(c, nil, "Server permanently deleted", http.StatusOK)
}

func (s *ServersController) Suspend(c *gin.Context) {
	id := c.Param("id")
	database.DB.Model(&models.Server{}).Where("id = ? OR uuid = ?", id, id).Update("suspended", "true")
	utils.Success(c, nil, "Server suspended", http.StatusOK)
}

func (s *ServersController) Unsuspend(c *gin.Context) {
	id := c.Param("id")
	database.DB.Model(&models.Server{}).Where("id = ? OR uuid = ?", id, id).Update("suspended", "false")
	utils.Success(c, nil, "Server unsuspended", http.StatusOK)
}

func (s *ServersController) GetByNode(c *gin.Context) {
	nodeID := c.Param("nodeId")
	var servers []models.Server
	database.DB.Where("node_id = ?", nodeID).Find(&servers)
	utils.Success(c, servers, "Servers retrieved", http.StatusOK)
}

func (s *ServersController) GetByOwner(c *gin.Context) {
	ownerID := c.Param("ownerId")
	var servers []models.Server
	database.DB.Where("owner_id = ?", ownerID).Find(&servers)
	utils.Success(c, servers, "Servers retrieved", http.StatusOK)
}
