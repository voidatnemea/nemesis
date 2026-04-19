package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mythicalsystems/nemesis-backend/internal/database"
	"github.com/mythicalsystems/nemesis-backend/internal/models"
	"github.com/mythicalsystems/nemesis-backend/pkg/utils"
)

type DatabasesController struct{}

func (d *DatabasesController) Index(c *gin.Context) {
	var dbs []models.DatabaseInstance
	database.DB.Find(&dbs)
	utils.Success(c, dbs, "Database instances retrieved", http.StatusOK)
}

func (d *DatabasesController) Show(c *gin.Context) {
	id := c.Param("id")
	var db models.DatabaseInstance
	if err := database.DB.First(&db, id).Error; err != nil {
		utils.Error(c, "Database instance not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	utils.Success(c, db, "Database instance retrieved", http.StatusOK)
}

func (d *DatabasesController) Create(c *gin.Context) {
	var req struct {
		Name              string `json:"name" binding:"required"`
		NodeID            int    `json:"node_id" binding:"required"`
		DatabaseType      string `json:"database_type"`
		DatabasePort      int    `json:"database_port"`
		DatabaseUsername  string `json:"database_username" binding:"required"`
		DatabasePassword  string `json:"database_password" binding:"required"`
		DatabaseHost      string `json:"database_host" binding:"required"`
		DatabaseSubdomain string `json:"database_subdomain"`
		MaxDatabases      int    `json:"max_databases"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	if req.DatabaseType == "" {
		req.DatabaseType = "mysql"
	}
	if req.DatabasePort == 0 {
		req.DatabasePort = 3306
	}
	db := models.DatabaseInstance{
		Name: req.Name, NodeID: req.NodeID, DatabaseType: req.DatabaseType,
		DatabasePort: req.DatabasePort, DatabaseUsername: req.DatabaseUsername,
		DatabasePassword: req.DatabasePassword, DatabaseHost: req.DatabaseHost,
		DatabaseSubdomain: req.DatabaseSubdomain, MaxDatabases: req.MaxDatabases,
	}
	database.DB.Create(&db)
	utils.Success(c, db, "Database instance created", http.StatusCreated)
}

func (d *DatabasesController) Update(c *gin.Context) {
	id := c.Param("id")
	var db models.DatabaseInstance
	if err := database.DB.First(&db, id).Error; err != nil {
		utils.Error(c, "Database instance not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	var req map[string]interface{}
	c.ShouldBindJSON(&req)
	delete(req, "id")
	database.DB.Model(&db).Updates(req)
	utils.Success(c, db, "Database instance updated", http.StatusOK)
}

func (d *DatabasesController) Delete(c *gin.Context) {
	id := c.Param("id")
	database.DB.Delete(&models.DatabaseInstance{}, id)
	utils.Success(c, nil, "Database instance deleted", http.StatusOK)
}

func (d *DatabasesController) GetByNode(c *gin.Context) {
	nodeID := c.Param("nodeId")
	var dbs []models.DatabaseInstance
	database.DB.Where("node_id = ?", nodeID).Find(&dbs)
	utils.Success(c, dbs, "Database instances retrieved", http.StatusOK)
}
