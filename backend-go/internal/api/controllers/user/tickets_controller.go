package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mythicalsystems/nemesis-backend/internal/database"
	"github.com/mythicalsystems/nemesis-backend/internal/models"
	"github.com/mythicalsystems/nemesis-backend/pkg/utils"
)

type UserTicketsController struct{}

func (t *UserTicketsController) Index(c *gin.Context) {
	ctxUser, _ := c.Get("user")
	user := ctxUser.(models.User)
	var tickets []models.Ticket
	database.DB.Where("user_uuid = ?", user.UUID).Find(&tickets)
	utils.Success(c, tickets, "Tickets retrieved", http.StatusOK)
}

func (t *UserTicketsController) Show(c *gin.Context) {
	ctxUser, _ := c.Get("user")
	user := ctxUser.(models.User)
	uuid := c.Param("uuid")
	var ticket models.Ticket
	if err := database.DB.Where("uuid = ? AND user_uuid = ?", uuid, user.UUID).First(&ticket).Error; err != nil {
		utils.Error(c, "Ticket not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	var messages []models.TicketMessage
	database.DB.Where("ticket_id = ?", ticket.ID).Find(&messages)
	utils.Success(c, gin.H{"ticket": ticket, "messages": messages}, "Ticket retrieved", http.StatusOK)
}

func (t *UserTicketsController) Create(c *gin.Context) {
	ctxUser, _ := c.Get("user")
	user := ctxUser.(models.User)
	var req struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description" binding:"required"`
		CategoryID  *int   `json:"category_id"`
		PriorityID  *int   `json:"priority_id"`
		ServerID    *int   `json:"server_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	ticket := models.Ticket{
		UUID:        utils.GenerateUUID(),
		Title:       req.Title,
		Description: req.Description,
		UserUUID:    user.UUID,
		CategoryID:  req.CategoryID,
		PriorityID:  req.PriorityID,
		ServerID:    req.ServerID,
	}
	database.DB.Create(&ticket)
	utils.Success(c, ticket, "Ticket created", http.StatusCreated)
}

func (t *UserTicketsController) Reply(c *gin.Context) {
	ctxUser, _ := c.Get("user")
	user := ctxUser.(models.User)
	uuid := c.Param("uuid")
	var ticket models.Ticket
	if err := database.DB.Where("uuid = ? AND user_uuid = ?", uuid, user.UUID).First(&ticket).Error; err != nil {
		utils.Error(c, "Ticket not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	var req struct {
		Message string `json:"message" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	msg := models.TicketMessage{
		TicketID: int(ticket.ID),
		UserUUID: &user.UUID,
		Message:  req.Message,
	}
	database.DB.Create(&msg)
	utils.Success(c, msg, "Reply sent", http.StatusCreated)
}

func (t *UserTicketsController) Delete(c *gin.Context) {
	ctxUser, _ := c.Get("user")
	user := ctxUser.(models.User)
	uuid := c.Param("uuid")
	database.DB.Where("uuid = ? AND user_uuid = ?", uuid, user.UUID).Delete(&models.Ticket{})
	utils.Success(c, nil, "Ticket deleted", http.StatusOK)
}

func (t *UserTicketsController) GetCategories(c *gin.Context) {
	var cats []models.TicketCategory
	database.DB.Find(&cats)
	utils.Success(c, cats, "Categories retrieved", http.StatusOK)
}

func (t *UserTicketsController) GetPriorities(c *gin.Context) {
	var priorities []models.TicketPriority
	database.DB.Find(&priorities)
	utils.Success(c, priorities, "Priorities retrieved", http.StatusOK)
}

func (t *UserTicketsController) GetStatuses(c *gin.Context) {
	var statuses []models.TicketStatus
	database.DB.Find(&statuses)
	utils.Success(c, statuses, "Statuses retrieved", http.StatusOK)
}
