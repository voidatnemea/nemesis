package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mythicalsystems/nemesis-backend/internal/database"
	"github.com/mythicalsystems/nemesis-backend/internal/models"
	"github.com/mythicalsystems/nemesis-backend/pkg/utils"
)

type NotificationsController struct{}

func (n *NotificationsController) Index(c *gin.Context) {
	p := utils.GetPagination(c)
	var notifications []models.Notification
	var total int64
	database.DB.Model(&models.Notification{}).Count(&total)
	database.DB.Order("created_at DESC").Offset(p.Offset).Limit(p.PerPage).Find(&notifications)
	utils.Success(c, gin.H{"notifications": notifications, "pagination": utils.BuildPagination(p, total)}, "Notifications retrieved", http.StatusOK)
}

func (n *NotificationsController) Show(c *gin.Context) {
	id := c.Param("id")
	var notif models.Notification
	if err := database.DB.First(&notif, id).Error; err != nil {
		utils.Error(c, "Notification not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	utils.Success(c, notif, "Notification retrieved", http.StatusOK)
}

func (n *NotificationsController) Create(c *gin.Context) {
	var req struct {
		Title           string `json:"title" binding:"required"`
		MessageMarkdown string `json:"message_markdown" binding:"required"`
		Type            string `json:"type" binding:"required"`
		IsDismissible   *bool  `json:"is_dismissible"`
		IsSticky        *bool  `json:"is_sticky"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	dismissible := "true"
	if req.IsDismissible != nil && !*req.IsDismissible {
		dismissible = "false"
	}
	sticky := "false"
	if req.IsSticky != nil && *req.IsSticky {
		sticky = "true"
	}
	notif := models.Notification{
		Title: req.Title, MessageMarkdown: req.MessageMarkdown, Type: req.Type,
		IsDismissible: dismissible, IsSticky: sticky,
	}
	database.DB.Create(&notif)
	utils.Success(c, notif, "Notification created", http.StatusCreated)
}

func (n *NotificationsController) Update(c *gin.Context) {
	id := c.Param("id")
	var notif models.Notification
	if err := database.DB.First(&notif, id).Error; err != nil {
		utils.Error(c, "Notification not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	var req map[string]interface{}
	c.ShouldBindJSON(&req)
	delete(req, "id")
	database.DB.Model(&notif).Updates(req)
	utils.Success(c, notif, "Notification updated", http.StatusOK)
}

func (n *NotificationsController) Delete(c *gin.Context) {
	id := c.Param("id")
	database.DB.Delete(&models.Notification{}, id)
	utils.Success(c, nil, "Notification deleted", http.StatusOK)
}
