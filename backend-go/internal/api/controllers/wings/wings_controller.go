package wings

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mythicalsystems/nemesis-backend/internal/database"
	"github.com/mythicalsystems/nemesis-backend/internal/models"
	"github.com/mythicalsystems/nemesis-backend/pkg/utils"
)

// WingsMiddleware authenticates requests coming from Wings nodes via daemon_token.
func WingsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		nodeID := c.Param("id")
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			utils.Error(c, "Unauthorized", "UNAUTHORIZED", http.StatusUnauthorized, nil)
			c.Abort()
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		var node models.Node
		if err := database.DB.Where("id = ? AND daemon_token = ?", nodeID, token).First(&node).Error; err != nil {
			utils.Error(c, "Invalid node token", "INVALID_NODE_TOKEN", http.StatusUnauthorized, nil)
			c.Abort()
			return
		}
		c.Set("node", node)
		c.Next()
	}
}

type WingsAdminController struct{}

func (w *WingsAdminController) Index(c *gin.Context) {
	var nodes []models.Node
	database.DB.Find(&nodes)
	utils.Success(c, nodes, "Nodes retrieved", http.StatusOK)
}

func (w *WingsAdminController) GetNodeInfo(c *gin.Context) {
	id := c.Param("id")
	var node models.Node
	if err := database.DB.First(&node, id).Error; err != nil {
		utils.Error(c, "Node not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	utils.Success(c, gin.H{
		"node":  node,
		"info":  "System information would be fetched from Wings daemon",
	}, "Node info retrieved", http.StatusOK)
}

func (w *WingsAdminController) GetServers(c *gin.Context) {
	id := c.Param("id")
	var servers []models.Server
	database.DB.Where("node_id = ?", id).Find(&servers)
	utils.Success(c, servers, "Node servers retrieved", http.StatusOK)
}
