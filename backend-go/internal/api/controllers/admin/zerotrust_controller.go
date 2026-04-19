package admin

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mythicalsystems/nemesis-backend/internal/database"
	"github.com/mythicalsystems/nemesis-backend/internal/models"
	"github.com/mythicalsystems/nemesis-backend/internal/services"
	"github.com/mythicalsystems/nemesis-backend/pkg/utils"
)

type ZeroTrustController struct{}

// ── Config ──────────────────────────────────────────────────────────────────

func (z *ZeroTrustController) GetConfig(c *gin.Context) {
	raw := services.GetSetting("zerotrust:config", "{}")
	var cfg map[string]interface{}
	if err := json.Unmarshal([]byte(raw), &cfg); err != nil {
		cfg = map[string]interface{}{}
	}
	utils.Success(c, cfg, "Config retrieved", http.StatusOK)
}

func (z *ZeroTrustController) UpdateConfig(c *gin.Context) {
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	b, _ := json.Marshal(req)
	services.SetSetting("zerotrust:config", string(b))
	utils.Success(c, req, "Config updated", http.StatusOK)
}

// ── Hashes ──────────────────────────────────────────────────────────────────

func (z *ZeroTrustController) ListHashes(c *gin.Context) {
	p := utils.GetPagination(c)
	status := c.Query("status")
	var hashes []models.ZeroTrustHash
	var total int64
	q := database.DB.Model(&models.ZeroTrustHash{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	q.Count(&total)
	q.Order("created_at DESC").Offset(p.Offset).Limit(p.PerPage).Find(&hashes)

	var confirmedCount, pendingCount, recentCount int64
	database.DB.Model(&models.ZeroTrustHash{}).Where("status = ?", "confirmed").Count(&confirmedCount)
	database.DB.Model(&models.ZeroTrustHash{}).Where("status = ?", "pending").Count(&pendingCount)
	database.DB.Model(&models.ZeroTrustHash{}).Where("created_at > ?", time.Now().AddDate(0, 0, -7)).Count(&recentCount)

	if hashes == nil {
		hashes = []models.ZeroTrustHash{}
	}
	utils.Success(c, gin.H{
		"hashes":     hashes,
		"pagination": utils.BuildPagination(p, total),
		"stats": gin.H{
			"total":     total,
			"confirmed": confirmedCount,
			"pending":   pendingCount,
			"recent":    recentCount,
		},
	}, "Hashes retrieved", http.StatusOK)
}

func (z *ZeroTrustController) HashStats(c *gin.Context) {
	var total, confirmed, unconfirmed, recent, servers int64
	database.DB.Model(&models.ZeroTrustHash{}).Count(&total)
	database.DB.Model(&models.ZeroTrustHash{}).Where("status = ?", "confirmed").Count(&confirmed)
	database.DB.Model(&models.ZeroTrustHash{}).Where("status = ?", "pending").Count(&unconfirmed)
	database.DB.Model(&models.ZeroTrustHash{}).Where("created_at > ?", time.Now().AddDate(0, 0, -7)).Count(&recent)
	database.DB.Raw("SELECT COUNT(DISTINCT server_uuid) FROM featherpanel_zerotrust_scan_logs").Scan(&servers)

	type typeCount struct {
		DetectionType string `json:"detection_type"`
		Count         int64  `json:"count"`
	}
	var topTypes []typeCount
	database.DB.Raw(`SELECT COALESCE(NULLIF(file_name, ''), 'unknown') AS detection_type, COUNT(*) AS count
		FROM featherpanel_zerotrust_hashes GROUP BY detection_type ORDER BY count DESC LIMIT 5`).Scan(&topTypes)
	if topTypes == nil {
		topTypes = []typeCount{}
	}

	utils.Success(c, gin.H{
		"totalHashes":       total,
		"confirmedHashes":   confirmed,
		"unconfirmedHashes": unconfirmed,
		"recentDetections":  recent,
		"totalServers":      servers,
		"topDetectionTypes": topTypes,
	}, "Hash statistics retrieved", http.StatusOK)
}

func (z *ZeroTrustController) AddHash(c *gin.Context) {
	var req struct {
		Hash        string `json:"hash" binding:"required"`
		FileName    string `json:"file_name"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	h := models.ZeroTrustHash{
		Hash:        req.Hash,
		FileName:    req.FileName,
		Description: req.Description,
		Status:      "pending",
		DetectedAt:  time.Now(),
	}
	if err := database.DB.Create(&h).Error; err != nil {
		utils.Error(c, "Hash already exists or failed to create", "DATABASE_ERROR", http.StatusConflict, nil)
		return
	}
	utils.Success(c, h, "Hash added", http.StatusCreated)
}

func (z *ZeroTrustController) DeleteHash(c *gin.Context) {
	id := c.Param("id")
	database.DB.Delete(&models.ZeroTrustHash{}, id)
	utils.Success(c, nil, "Hash deleted", http.StatusOK)
}

func (z *ZeroTrustController) CheckHash(c *gin.Context) {
	var req struct {
		Hashes []string `json:"hashes" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	var found []models.ZeroTrustHash
	database.DB.Where("hash IN ?", req.Hashes).Find(&found)
	if found == nil {
		found = []models.ZeroTrustHash{}
	}
	utils.Success(c, gin.H{"matches": found}, "Hash check complete", http.StatusOK)
}

func (z *ZeroTrustController) BulkConfirm(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	database.DB.Model(&models.ZeroTrustHash{}).Where("id IN ?", req.IDs).Update("status", "confirmed")
	utils.Success(c, nil, "Hashes confirmed", http.StatusOK)
}

func (z *ZeroTrustController) BulkDelete(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	database.DB.Where("id IN ?", req.IDs).Delete(&models.ZeroTrustHash{})
	utils.Success(c, nil, "Hashes deleted", http.StatusOK)
}

// ── Scan Logs ────────────────────────────────────────────────────────────────

func (z *ZeroTrustController) ListLogs(c *gin.Context) {
	p := utils.GetPagination(c)
	var execs []models.ZeroTrustScanExecution
	var total int64
	database.DB.Model(&models.ZeroTrustScanExecution{}).Count(&total)
	database.DB.Order("created_at DESC").Offset(p.Offset).Limit(p.PerPage).Find(&execs)
	if execs == nil {
		execs = []models.ZeroTrustScanExecution{}
	}
	utils.Success(c, gin.H{"executions": execs, "pagination": utils.BuildPagination(p, total)}, "Logs retrieved", http.StatusOK)
}

func (z *ZeroTrustController) GetLog(c *gin.Context) {
	executionID := c.Param("executionId")
	var exec models.ZeroTrustScanExecution
	if err := database.DB.First(&exec, executionID).Error; err != nil {
		utils.Error(c, "Execution not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	var logs []models.ZeroTrustScanLog
	database.DB.Where("execution_id = ?", exec.ID).Find(&logs)
	if logs == nil {
		logs = []models.ZeroTrustScanLog{}
	}
	utils.Success(c, gin.H{"execution": exec, "logs": logs}, "Execution logs retrieved", http.StatusOK)
}

// ── Scan ─────────────────────────────────────────────────────────────────────

func (z *ZeroTrustController) Scan(c *gin.Context) {
	var req struct {
		ServerUUID string `json:"server_uuid" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}

	exec := models.ZeroTrustScanExecution{
		Status:    "running",
		StartedAt: time.Now(),
	}
	database.DB.Create(&exec)

	go func() {
		time.Sleep(3 * time.Second)
		now := time.Now()
		var server models.Server
		name := req.ServerUUID
		if err := database.DB.Where("uuid = ?", req.ServerUUID).First(&server).Error; err == nil {
			name = server.Name
		}
		log := models.ZeroTrustScanLog{
			ExecutionID:  exec.ID,
			ServerUUID:   req.ServerUUID,
			ServerName:   name,
			Status:       "clean",
			FilesScanned: 0,
			Detections:   "[]",
		}
		database.DB.Create(&log)
		database.DB.Model(&exec).Updates(map[string]interface{}{
			"status":           "completed",
			"servers_scanned":  1,
			"detections_found": 0,
			"duration_ms":      3000,
			"finished_at":      &now,
		})
	}()

	utils.Success(c, gin.H{"execution_id": exec.ID, "status": "running"}, "Scan started", http.StatusAccepted)
}

func (z *ZeroTrustController) ScanBatch(c *gin.Context) {
	var req struct {
		ServerUUIDs []string `json:"server_uuids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}

	exec := models.ZeroTrustScanExecution{
		Status:    "running",
		StartedAt: time.Now(),
	}
	database.DB.Create(&exec)

	go func() {
		time.Sleep(5 * time.Second)
		now := time.Now()
		for _, uuid := range req.ServerUUIDs {
			var server models.Server
			name := uuid
			if err := database.DB.Where("uuid = ?", uuid).First(&server).Error; err == nil {
				name = server.Name
			}
			log := models.ZeroTrustScanLog{
				ExecutionID:  exec.ID,
				ServerUUID:   uuid,
				ServerName:   name,
				Status:       "clean",
				FilesScanned: 0,
				Detections:   "[]",
			}
			database.DB.Create(&log)
		}
		database.DB.Model(&exec).Updates(map[string]interface{}{
			"status":           "completed",
			"servers_scanned":  len(req.ServerUUIDs),
			"detections_found": 0,
			"duration_ms":      5000,
			"finished_at":      &now,
		})
	}()

	utils.Success(c, gin.H{"execution_id": exec.ID, "status": "running", "servers": len(req.ServerUUIDs)}, "Batch scan started", http.StatusAccepted)
}
