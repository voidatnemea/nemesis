package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mythicalsystems/nemesis-backend/internal/database"
	"github.com/mythicalsystems/nemesis-backend/internal/models"
	"github.com/mythicalsystems/nemesis-backend/pkg/utils"
)

type NodesController struct{}

func (n *NodesController) Index(c *gin.Context) {
	p := utils.GetPagination(c)
	search := c.Query("search")
	var nodes []models.Node
	var total int64
	q := database.DB.Model(&models.Node{})
	if search != "" {
		q = q.Where("name LIKE ? OR fqdn LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	q.Count(&total)
	q.Offset(p.Offset).Limit(p.PerPage).Find(&nodes)
	utils.Success(c, gin.H{"nodes": nodes, "pagination": utils.BuildPagination(p, total)}, "Nodes retrieved", http.StatusOK)
}

func (n *NodesController) GetIPs(c *gin.Context) {
	nodeID := c.Param("id")
	type ipRow struct {
		IP string
	}
	var rows []ipRow
	database.DB.Raw("SELECT DISTINCT ip FROM featherpanel_allocations WHERE node_id = ?", nodeID).Scan(&rows)
	ips := make([]string, 0, len(rows))
	for _, r := range rows {
		ips = append(ips, r.IP)
	}
	utils.Success(c, gin.H{"ips": gin.H{"ip_addresses": ips}}, "IPs retrieved", http.StatusOK)
}

func (n *NodesController) Show(c *gin.Context) {
	id := c.Param("id")
	var node models.Node
	if err := database.DB.First(&node, id).Error; err != nil {
		utils.Error(c, "Node not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	utils.Success(c, node, "Node retrieved", http.StatusOK)
}

func (n *NodesController) Create(c *gin.Context) {
	var req struct {
		Name            string `json:"name" binding:"required"`
		Description     string `json:"description"`
		LocationID      int    `json:"location_id" binding:"required"`
		FQDN            string `json:"fqdn" binding:"required"`
		Public          string `json:"public"`
		Scheme          string `json:"scheme"`
		BehindProxy     string `json:"behind_proxy"`
		Memory          int64  `json:"memory" binding:"required"`
		MemoryOvercommit int   `json:"memory_overcommit"`
		Disk            int64  `json:"disk" binding:"required"`
		DiskOvercommit  int    `json:"disk_overcommit"`
		DaemonBase      string `json:"daemon_base"`
		DaemonSFTP      int    `json:"daemon_sftp"`
		DaemonListen    int    `json:"daemon_listen"`
		UploadSize      int    `json:"upload_size"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	if req.Scheme == "" {
		req.Scheme = "https"
	}
	if req.DaemonSFTP == 0 {
		req.DaemonSFTP = 2022
	}
	if req.DaemonListen == 0 {
		req.DaemonListen = 8080
	}
	if req.UploadSize == 0 {
		req.UploadSize = 100
	}
	node := models.Node{
		UUID:             utils.GenerateUUID(),
		Name:             req.Name,
		Description:      req.Description,
		LocationID:       req.LocationID,
		FQDN:             req.FQDN,
		Public:           req.Public,
		Scheme:           req.Scheme,
		BehindProxy:      req.BehindProxy,
		Memory:           req.Memory,
		MemoryOvercommit: req.MemoryOvercommit,
		Disk:             req.Disk,
		DiskOvercommit:   req.DiskOvercommit,
		DaemonBase:       req.DaemonBase,
		DaemonSFTP:       req.DaemonSFTP,
		DaemonListen:     req.DaemonListen,
		DaemonToken:      utils.GenerateRandomString(32),
		UploadSize:       req.UploadSize,
	}
	if err := database.DB.Create(&node).Error; err != nil {
		utils.Error(c, "Failed to create node", "DATABASE_ERROR", http.StatusInternalServerError, nil)
		return
	}
	utils.Success(c, node, "Node created", http.StatusCreated)
}

func (n *NodesController) Update(c *gin.Context) {
	id := c.Param("id")
	var node models.Node
	if err := database.DB.First(&node, id).Error; err != nil {
		utils.Error(c, "Node not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	delete(req, "id")
	delete(req, "uuid")
	delete(req, "daemon_token")
	database.DB.Model(&node).Updates(req)
	utils.Success(c, node, "Node updated", http.StatusOK)
}

func (n *NodesController) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.Node{}, id).Error; err != nil {
		utils.Error(c, "Failed to delete node", "DATABASE_ERROR", http.StatusInternalServerError, nil)
		return
	}
	utils.Success(c, nil, "Node deleted", http.StatusOK)
}

func (n *NodesController) ResetKey(c *gin.Context) {
	id := c.Param("id")
	newToken := utils.GenerateRandomString(32)
	if err := database.DB.Model(&models.Node{}).Where("id = ?", id).Update("daemon_token", newToken).Error; err != nil {
		utils.Error(c, "Failed to reset key", "DATABASE_ERROR", http.StatusInternalServerError, nil)
		return
	}
	utils.Success(c, gin.H{"token": newToken}, "Node key reset", http.StatusOK)
}
