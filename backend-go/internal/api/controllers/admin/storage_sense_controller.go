package admin

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mythicalsystems/nemesis-backend/internal/database"
	"github.com/mythicalsystems/nemesis-backend/pkg/utils"
)

type StorageSenseController struct{}

type storageCategory struct {
	ID               string `json:"id"`
	Table            string `json:"table"`
	Available        bool   `json:"available"`
	UsesRetentionDays bool   `json:"uses_retention_days"`
	RowCount         int64  `json:"row_count"`
	PurgeableCount   int64  `json:"purgeable_count"`
	ApproxDataBytes  int64  `json:"approx_data_bytes"`
}

var storageCategoryDefs = []struct {
	ID    string
	Table string
	Has   bool
}{
	{"user_activity", "featherpanel_activity", true},
	{"server_activity", "featherpanel_server_activities", true},
	{"mail_history", "featherpanel_mail_queues", true},
	{"admin_notifications", "featherpanel_notifications", true},
	{"sso_expired_tokens", "featherpanel_sso_tokens", true},
	{"vm_instance_activity", "featherpanel_vm_activities", false},
	{"vm_panel_logs", "featherpanel_vm_panel_logs", false},
	{"chatbot_data", "featherpanel_chatbot_data", false},
	{"featherzerotrust_logs", "featherpanel_zerotrust_scan_logs", false},
}

func (s *StorageSenseController) Summary(c *gin.Context) {
	daysStr := c.DefaultQuery("days", "90")
	days, _ := strconv.Atoi(daysStr)
	if days < 7 {
		days = 7
	}

	cutoff := time.Now().AddDate(0, 0, -days)

	var categories []storageCategory
	var totalRows, totalPurgeable, totalBytes int64

	for _, def := range storageCategoryDefs {
		cat := storageCategory{
			ID:                def.ID,
			Table:             def.Table,
			Available:         def.Has,
			UsesRetentionDays: true,
		}

		if def.Has {
			db := database.DB.Raw("SELECT COUNT(*) FROM `"+def.Table+"`").Scan(&cat.RowCount)
			if db.Error != nil {
				cat.RowCount = 0
				cat.Available = false
			}
			database.DB.Raw("SELECT COUNT(*) FROM `"+def.Table+"` WHERE created_at < ?", cutoff).Scan(&cat.PurgeableCount)
			cat.ApproxDataBytes = cat.RowCount * 512
			totalRows += cat.RowCount
			totalPurgeable += cat.PurgeableCount
			totalBytes += cat.ApproxDataBytes
		}

		categories = append(categories, cat)
	}

	utils.Success(c, gin.H{
		"categories": categories,
		"totals": gin.H{
			"tables_tracked":   len(storageCategoryDefs),
			"total_rows":       totalRows,
			"total_purgeable":  totalPurgeable,
			"approx_data_bytes": totalBytes,
		},
		"disk": gin.H{
			"path":  "/opt/nemesis",
			"bytes": int64(0),
		},
	}, "Storage sense data retrieved", http.StatusOK)
}

func (s *StorageSenseController) Purge(c *gin.Context) {
	var req struct {
		Category string `json:"category" binding:"required"`
		Days     int    `json:"days"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	if req.Days < 7 {
		req.Days = 90
	}
	cutoff := time.Now().AddDate(0, 0, -req.Days)

	for _, def := range storageCategoryDefs {
		if def.ID == req.Category && def.Has {
			res := database.DB.Exec("DELETE FROM `"+def.Table+"` WHERE created_at < ?", cutoff)
			utils.Success(c, gin.H{"deleted": res.RowsAffected}, "Category purged", http.StatusOK)
			return
		}
	}
	utils.Error(c, "Category not found or not available", "NOT_FOUND", http.StatusNotFound, nil)
}

func (s *StorageSenseController) PurgeBatch(c *gin.Context) {
	var req struct {
		Categories []string `json:"categories" binding:"required"`
		Days       int      `json:"days"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	if req.Days < 7 {
		req.Days = 90
	}
	cutoff := time.Now().AddDate(0, 0, -req.Days)

	catSet := make(map[string]bool)
	for _, id := range req.Categories {
		catSet[id] = true
	}

	var totalDeleted int64
	for _, def := range storageCategoryDefs {
		if catSet[def.ID] && def.Has {
			res := database.DB.Exec("DELETE FROM `"+def.Table+"` WHERE created_at < ?", cutoff)
			totalDeleted += res.RowsAffected
		}
	}

	utils.Success(c, gin.H{"deleted": totalDeleted}, "Batch purge completed", http.StatusOK)
}
