package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mythicalsystems/nemesis-backend/internal/database"
	"github.com/mythicalsystems/nemesis-backend/internal/models"
	"github.com/mythicalsystems/nemesis-backend/pkg/utils"
)

type ServerUserController struct{}

func (s *ServerUserController) GetUserServers(c *gin.Context) {
	ctxUser, _ := c.Get("user")
	user := ctxUser.(models.User)
	var servers []models.Server
	database.DB.Where("owner_id = ?", user.ID).Find(&servers)
	utils.Success(c, servers, "Servers retrieved", http.StatusOK)
}

func (s *ServerUserController) GetServer(c *gin.Context) {
	ctxUser, _ := c.Get("user")
	user := ctxUser.(models.User)
	uuidShort := c.Param("uuidShort")
	var server models.Server
	if err := database.DB.Where("(uuid_short = ? OR uuid = ?) AND (owner_id = ? OR id IN (SELECT server_id FROM featherpanel_server_subusers WHERE user_id = ?))",
		uuidShort, uuidShort, user.ID, user.ID).First(&server).Error; err != nil {
		utils.Error(c, "Server not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	utils.Success(c, server, "Server retrieved", http.StatusOK)
}

func (s *ServerUserController) UpdateServer(c *gin.Context) {
	ctxUser, _ := c.Get("user")
	user := ctxUser.(models.User)
	uuidShort := c.Param("uuidShort")
	var server models.Server
	if err := database.DB.Where("(uuid_short = ? OR uuid = ?) AND owner_id = ?", uuidShort, uuidShort, user.ID).First(&server).Error; err != nil {
		utils.Error(c, "Server not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	c.ShouldBindJSON(&req)
	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	database.DB.Model(&server).Updates(updates)
	utils.Success(c, server, "Server updated", http.StatusOK)
}

func (s *ServerUserController) DeleteServer(c *gin.Context) {
	ctxUser, _ := c.Get("user")
	user := ctxUser.(models.User)
	uuidShort := c.Param("uuidShort")
	database.DB.Where("(uuid_short = ? OR uuid = ?) AND owner_id = ?", uuidShort, uuidShort, user.ID).Delete(&models.Server{})
	utils.Success(c, nil, "Server deleted", http.StatusOK)
}

// ServerBackupController
type ServerBackupController struct{}

func (b *ServerBackupController) GetBackups(c *gin.Context) {
	server := getServerFromContext(c)
	if server == nil {
		return
	}
	var backups []models.Backup
	database.DB.Where("server_id = ?", server.ID).Order("created_at DESC").Find(&backups)
	utils.Success(c, backups, "Backups retrieved", http.StatusOK)
}

func (b *ServerBackupController) CreateBackup(c *gin.Context) {
	server := getServerFromContext(c)
	if server == nil {
		return
	}
	var req struct {
		Name         string `json:"name"`
		IgnoredFiles string `json:"ignored_files"`
	}
	c.ShouldBindJSON(&req)
	if req.Name == "" {
		req.Name = "Backup"
	}
	backup := models.Backup{
		ServerID:     int(server.ID),
		UUID:         utils.GenerateUUID(),
		Name:         req.Name,
		IgnoredFiles: req.IgnoredFiles,
		Disk:         "local",
	}
	database.DB.Create(&backup)
	utils.Success(c, backup, "Backup created", http.StatusCreated)
}

func (b *ServerBackupController) GetBackup(c *gin.Context) {
	server := getServerFromContext(c)
	if server == nil {
		return
	}
	backupUUID := c.Param("backupUuid")
	var backup models.Backup
	if err := database.DB.Where("uuid = ? AND server_id = ?", backupUUID, server.ID).First(&backup).Error; err != nil {
		utils.Error(c, "Backup not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	utils.Success(c, backup, "Backup retrieved", http.StatusOK)
}

func (b *ServerBackupController) DeleteBackup(c *gin.Context) {
	server := getServerFromContext(c)
	if server == nil {
		return
	}
	backupUUID := c.Param("backupUuid")
	database.DB.Where("uuid = ? AND server_id = ?", backupUUID, server.ID).Delete(&models.Backup{})
	utils.Success(c, nil, "Backup deleted", http.StatusOK)
}

func (b *ServerBackupController) LockBackup(c *gin.Context) {
	server := getServerFromContext(c)
	if server == nil {
		return
	}
	backupUUID := c.Param("backupUuid")
	database.DB.Model(&models.Backup{}).Where("uuid = ? AND server_id = ?", backupUUID, server.ID).Update("is_locked", "true")
	utils.Success(c, nil, "Backup locked", http.StatusOK)
}

func (b *ServerBackupController) UnlockBackup(c *gin.Context) {
	server := getServerFromContext(c)
	if server == nil {
		return
	}
	backupUUID := c.Param("backupUuid")
	database.DB.Model(&models.Backup{}).Where("uuid = ? AND server_id = ?", backupUUID, server.ID).Update("is_locked", "false")
	utils.Success(c, nil, "Backup unlocked", http.StatusOK)
}

// ServerDatabaseController
type ServerDatabaseController struct{}

func (d *ServerDatabaseController) GetServerDatabases(c *gin.Context) {
	server := getServerFromContext(c)
	if server == nil {
		return
	}
	var dbs []models.ServerDatabase
	database.DB.Where("server_id = ?", server.ID).Find(&dbs)
	utils.Success(c, dbs, "Databases retrieved", http.StatusOK)
}

func (d *ServerDatabaseController) CreateServerDatabase(c *gin.Context) {
	server := getServerFromContext(c)
	if server == nil {
		return
	}
	var req struct {
		DatabaseHostID int    `json:"database_host_id" binding:"required"`
		Database       string `json:"database" binding:"required"`
		Username       string `json:"username" binding:"required"`
		Password       string `json:"password" binding:"required"`
		Remote         string `json:"remote"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	if req.Remote == "" {
		req.Remote = "%"
	}
	db := models.ServerDatabase{
		ServerID: int(server.ID), DatabaseHostID: req.DatabaseHostID,
		Database: req.Database, Username: req.Username, Password: req.Password, Remote: req.Remote,
	}
	database.DB.Create(&db)
	utils.Success(c, db, "Database created", http.StatusCreated)
}

func (d *ServerDatabaseController) GetServerDatabase(c *gin.Context) {
	server := getServerFromContext(c)
	if server == nil {
		return
	}
	dbID := c.Param("databaseId")
	var db models.ServerDatabase
	if err := database.DB.Where("id = ? AND server_id = ?", dbID, server.ID).First(&db).Error; err != nil {
		utils.Error(c, "Database not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	utils.Success(c, db, "Database retrieved", http.StatusOK)
}

func (d *ServerDatabaseController) DeleteServerDatabase(c *gin.Context) {
	server := getServerFromContext(c)
	if server == nil {
		return
	}
	dbID := c.Param("databaseId")
	database.DB.Where("id = ? AND server_id = ?", dbID, server.ID).Delete(&models.ServerDatabase{})
	utils.Success(c, nil, "Database deleted", http.StatusOK)
}

func (d *ServerDatabaseController) GetAvailableDatabaseHosts(c *gin.Context) {
	var hosts []models.DatabaseInstance
	database.DB.Find(&hosts)
	utils.Success(c, hosts, "Available database hosts retrieved", http.StatusOK)
}

// ServerAllocationController
type ServerAllocationController struct{}

func (a *ServerAllocationController) GetServerAllocations(c *gin.Context) {
	server := getServerFromContext(c)
	if server == nil {
		return
	}
	var allocs []models.Allocation
	database.DB.Where("server_id = ?", server.ID).Find(&allocs)
	utils.Success(c, allocs, "Allocations retrieved", http.StatusOK)
}

func (a *ServerAllocationController) DeleteAllocation(c *gin.Context) {
	server := getServerFromContext(c)
	if server == nil {
		return
	}
	allocationID := c.Param("allocationId")
	// Can't delete primary allocation
	if allocationID == string(rune(server.AllocationID)) {
		utils.Error(c, "Cannot delete primary allocation", "CANNOT_DELETE_PRIMARY", http.StatusBadRequest, nil)
		return
	}
	database.DB.Model(&models.Allocation{}).Where("id = ? AND server_id = ?", allocationID, server.ID).Update("server_id", nil)
	utils.Success(c, nil, "Allocation removed", http.StatusOK)
}

func (a *ServerAllocationController) SetPrimaryAllocation(c *gin.Context) {
	server := getServerFromContext(c)
	if server == nil {
		return
	}
	allocationID := c.Param("allocationId")
	database.DB.Model(&server).Update("allocation_id", allocationID)
	utils.Success(c, nil, "Primary allocation updated", http.StatusOK)
}

// ServerScheduleController
type ServerScheduleController struct{}

func (s *ServerScheduleController) GetSchedules(c *gin.Context) {
	server := getServerFromContext(c)
	if server == nil {
		return
	}
	var schedules []models.ServerSchedule
	database.DB.Where("server_id = ?", server.ID).Find(&schedules)
	utils.Success(c, schedules, "Schedules retrieved", http.StatusOK)
}

func (s *ServerScheduleController) CreateSchedule(c *gin.Context) {
	server := getServerFromContext(c)
	if server == nil {
		return
	}
	var req struct {
		Name           string `json:"name" binding:"required"`
		CronDayOfWeek  string `json:"cron_day_of_week"`
		CronMonth      string `json:"cron_month"`
		CronDayOfMonth string `json:"cron_day_of_month"`
		CronHour       string `json:"cron_hour"`
		CronMinute     string `json:"cron_minute"`
		IsActive       string `json:"is_active"`
		OnlyWhenOnline string `json:"only_when_online"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	schedule := models.ServerSchedule{
		ServerID: int(server.ID), Name: req.Name,
		CronDayOfWeek: req.CronDayOfWeek, CronMonth: req.CronMonth,
		CronDayOfMonth: req.CronDayOfMonth, CronHour: req.CronHour, CronMinute: req.CronMinute,
		IsActive: req.IsActive, OnlyWhenOnline: req.OnlyWhenOnline,
	}
	if schedule.IsActive == "" {
		schedule.IsActive = "true"
	}
	database.DB.Create(&schedule)
	utils.Success(c, schedule, "Schedule created", http.StatusCreated)
}

func (s *ServerScheduleController) GetSchedule(c *gin.Context) {
	server := getServerFromContext(c)
	if server == nil {
		return
	}
	scheduleID := c.Param("scheduleId")
	var schedule models.ServerSchedule
	if err := database.DB.Where("id = ? AND server_id = ?", scheduleID, server.ID).First(&schedule).Error; err != nil {
		utils.Error(c, "Schedule not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	utils.Success(c, schedule, "Schedule retrieved", http.StatusOK)
}

func (s *ServerScheduleController) UpdateSchedule(c *gin.Context) {
	server := getServerFromContext(c)
	if server == nil {
		return
	}
	scheduleID := c.Param("scheduleId")
	var schedule models.ServerSchedule
	if err := database.DB.Where("id = ? AND server_id = ?", scheduleID, server.ID).First(&schedule).Error; err != nil {
		utils.Error(c, "Schedule not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	var req map[string]interface{}
	c.ShouldBindJSON(&req)
	delete(req, "id")
	database.DB.Model(&schedule).Updates(req)
	utils.Success(c, schedule, "Schedule updated", http.StatusOK)
}

func (s *ServerScheduleController) DeleteSchedule(c *gin.Context) {
	server := getServerFromContext(c)
	if server == nil {
		return
	}
	scheduleID := c.Param("scheduleId")
	database.DB.Where("id = ? AND server_id = ?", scheduleID, server.ID).Delete(&models.ServerSchedule{})
	utils.Success(c, nil, "Schedule deleted", http.StatusOK)
}

func (s *ServerScheduleController) ToggleSchedule(c *gin.Context) {
	server := getServerFromContext(c)
	if server == nil {
		return
	}
	scheduleID := c.Param("scheduleId")
	var schedule models.ServerSchedule
	if err := database.DB.Where("id = ? AND server_id = ?", scheduleID, server.ID).First(&schedule).Error; err != nil {
		utils.Error(c, "Schedule not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	newActive := "true"
	if schedule.IsActive == "true" {
		newActive = "false"
	}
	database.DB.Model(&schedule).Update("is_active", newActive)
	utils.Success(c, schedule, "Schedule toggled", http.StatusOK)
}

// helper — looks up a server by uuidShort and verifies caller has access
func getServerFromContext(c *gin.Context) *models.Server {
	ctxUser, _ := c.Get("user")
	user := ctxUser.(models.User)
	uuidShort := c.Param("uuidShort")
	var server models.Server
	if err := database.DB.Where(
		"(uuid_short = ? OR uuid = ?) AND (owner_id = ? OR id IN (SELECT server_id FROM featherpanel_server_subusers WHERE user_id = ?))",
		uuidShort, uuidShort, user.ID, user.ID,
	).First(&server).Error; err != nil {
		utils.Error(c, "Server not found", "NOT_FOUND", http.StatusNotFound, nil)
		c.Abort()
		return nil
	}
	return &server
}
