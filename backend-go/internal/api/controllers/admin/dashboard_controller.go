package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mythicalsystems/nemesis-backend/internal/database"
	"github.com/mythicalsystems/nemesis-backend/internal/models"
	"github.com/mythicalsystems/nemesis-backend/pkg/utils"
)

type DashboardController struct{}

func (d *DashboardController) Index(c *gin.Context) {
	var userCount, serverCount, nodeCount, ticketCount int64
	database.DB.Model(&models.User{}).Count(&userCount)
	database.DB.Model(&models.Server{}).Count(&serverCount)
	database.DB.Model(&models.Node{}).Count(&nodeCount)
	database.DB.Model(&models.Ticket{}).Count(&ticketCount)

	utils.Success(c, gin.H{
		"users":   userCount,
		"servers": serverCount,
		"nodes":   nodeCount,
		"tickets": ticketCount,
	}, "Dashboard data retrieved", http.StatusOK)
}
