package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mythicalsystems/nemesis-backend/internal/database"
	"github.com/mythicalsystems/nemesis-backend/internal/models"
	"github.com/mythicalsystems/nemesis-backend/pkg/utils"
)

type AdminTicketsController struct{}

func (t *AdminTicketsController) Index(c *gin.Context) {
	p := utils.GetPagination(c)
	search := c.Query("search")
	var tickets []models.Ticket
	var total int64
	q := database.DB.Model(&models.Ticket{})
	if search != "" {
		q = q.Where("title LIKE ? OR description LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	q.Count(&total)
	q.Offset(p.Offset).Limit(p.PerPage).Find(&tickets)
	utils.Success(c, gin.H{"tickets": tickets, "pagination": utils.BuildPagination(p, total)}, "Tickets retrieved", http.StatusOK)
}

func (t *AdminTicketsController) Show(c *gin.Context) {
	uuid := c.Param("uuid")
	var ticket models.Ticket
	if err := database.DB.Where("uuid = ?", uuid).First(&ticket).Error; err != nil {
		utils.Error(c, "Ticket not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	utils.Success(c, ticket, "Ticket retrieved", http.StatusOK)
}

func (t *AdminTicketsController) Update(c *gin.Context) {
	uuid := c.Param("uuid")
	var ticket models.Ticket
	if err := database.DB.Where("uuid = ?", uuid).First(&ticket).Error; err != nil {
		utils.Error(c, "Ticket not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	var req map[string]interface{}
	c.ShouldBindJSON(&req)
	delete(req, "id")
	delete(req, "uuid")
	database.DB.Model(&ticket).Updates(req)
	utils.Success(c, ticket, "Ticket updated", http.StatusOK)
}

func (t *AdminTicketsController) Delete(c *gin.Context) {
	uuid := c.Param("uuid")
	database.DB.Where("uuid = ?", uuid).Delete(&models.Ticket{})
	utils.Success(c, nil, "Ticket deleted", http.StatusOK)
}

func (t *AdminTicketsController) Reply(c *gin.Context) {
	uuid := c.Param("uuid")
	var ticket models.Ticket
	if err := database.DB.Where("uuid = ?", uuid).First(&ticket).Error; err != nil {
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
	ctxUser, _ := c.Get("user")
	user := ctxUser.(models.User)
	msg := models.TicketMessage{
		TicketID: int(ticket.ID),
		UserUUID: &user.UUID,
		Message:  req.Message,
		IsStaff:  "true",
	}
	database.DB.Create(&msg)
	utils.Success(c, msg, "Reply sent", http.StatusCreated)
}

func (t *AdminTicketsController) Close(c *gin.Context) {
	uuid := c.Param("uuid")
	// Find a "closed" status
	var status models.TicketStatus
	database.DB.Where("name LIKE ?", "%closed%").First(&status)
	database.DB.Model(&models.Ticket{}).Where("uuid = ?", uuid).Update("status_id", status.ID)
	utils.Success(c, nil, "Ticket closed", http.StatusOK)
}

func (t *AdminTicketsController) Reopen(c *gin.Context) {
	uuid := c.Param("uuid")
	var status models.TicketStatus
	database.DB.Where("name LIKE ?", "%open%").First(&status)
	database.DB.Model(&models.Ticket{}).Where("uuid = ?", uuid).Update("status_id", status.ID)
	utils.Success(c, nil, "Ticket reopened", http.StatusOK)
}

// TicketCategoriesController
type TicketCategoriesController struct{}

func (tc *TicketCategoriesController) Index(c *gin.Context) {
	var cats []models.TicketCategory
	database.DB.Find(&cats)
	utils.Success(c, gin.H{"categories": cats}, "Categories retrieved", http.StatusOK)
}

func (tc *TicketCategoriesController) Show(c *gin.Context) {
	id := c.Param("id")
	var cat models.TicketCategory
	if err := database.DB.First(&cat, id).Error; err != nil {
		utils.Error(c, "Category not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	utils.Success(c, cat, "Category retrieved", http.StatusOK)
}

func (tc *TicketCategoriesController) Create(c *gin.Context) {
	var req struct {
		Name         string `json:"name" binding:"required"`
		Icon         string `json:"icon"`
		Color        string `json:"color"`
		SupportEmail string `json:"support_email"`
		OpenHours    string `json:"open_hours"`
		Description  string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	cat := models.TicketCategory{
		Name: req.Name, Icon: req.Icon, Color: req.Color,
		SupportEmail: req.SupportEmail, OpenHours: req.OpenHours, Description: req.Description,
	}
	database.DB.Create(&cat)
	utils.Success(c, cat, "Category created", http.StatusCreated)
}

func (tc *TicketCategoriesController) Update(c *gin.Context) {
	id := c.Param("id")
	var cat models.TicketCategory
	if err := database.DB.First(&cat, id).Error; err != nil {
		utils.Error(c, "Category not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	var req map[string]interface{}
	c.ShouldBindJSON(&req)
	delete(req, "id")
	database.DB.Model(&cat).Updates(req)
	utils.Success(c, cat, "Category updated", http.StatusOK)
}

func (tc *TicketCategoriesController) Delete(c *gin.Context) {
	id := c.Param("id")
	database.DB.Delete(&models.TicketCategory{}, id)
	utils.Success(c, nil, "Category deleted", http.StatusOK)
}

// TicketStatusesController
type TicketStatusesController struct{}

func (ts *TicketStatusesController) Index(c *gin.Context) {
	var statuses []models.TicketStatus
	database.DB.Find(&statuses)
	utils.Success(c, gin.H{"statuses": statuses}, "Statuses retrieved", http.StatusOK)
}

func (ts *TicketStatusesController) Show(c *gin.Context) {
	id := c.Param("id")
	var status models.TicketStatus
	if err := database.DB.First(&status, id).Error; err != nil {
		utils.Error(c, "Status not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	utils.Success(c, status, "Status retrieved", http.StatusOK)
}

func (ts *TicketStatusesController) Create(c *gin.Context) {
	var req struct {
		Name  string `json:"name" binding:"required"`
		Color string `json:"color"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	status := models.TicketStatus{Name: req.Name, Color: req.Color}
	database.DB.Create(&status)
	utils.Success(c, status, "Status created", http.StatusCreated)
}

func (ts *TicketStatusesController) Update(c *gin.Context) {
	id := c.Param("id")
	var status models.TicketStatus
	if err := database.DB.First(&status, id).Error; err != nil {
		utils.Error(c, "Status not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	var req map[string]interface{}
	c.ShouldBindJSON(&req)
	delete(req, "id")
	database.DB.Model(&status).Updates(req)
	utils.Success(c, status, "Status updated", http.StatusOK)
}

func (ts *TicketStatusesController) Delete(c *gin.Context) {
	id := c.Param("id")
	database.DB.Delete(&models.TicketStatus{}, id)
	utils.Success(c, nil, "Status deleted", http.StatusOK)
}

// TicketPrioritiesController
type TicketPrioritiesController struct{}

func (tp *TicketPrioritiesController) Index(c *gin.Context) {
	var priorities []models.TicketPriority
	database.DB.Find(&priorities)
	utils.Success(c, gin.H{"priorities": priorities}, "Priorities retrieved", http.StatusOK)
}

func (tp *TicketPrioritiesController) Show(c *gin.Context) {
	id := c.Param("id")
	var priority models.TicketPriority
	if err := database.DB.First(&priority, id).Error; err != nil {
		utils.Error(c, "Priority not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	utils.Success(c, priority, "Priority retrieved", http.StatusOK)
}

func (tp *TicketPrioritiesController) Create(c *gin.Context) {
	var req struct {
		Name  string `json:"name" binding:"required"`
		Color string `json:"color"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	priority := models.TicketPriority{Name: req.Name, Color: req.Color}
	database.DB.Create(&priority)
	utils.Success(c, priority, "Priority created", http.StatusCreated)
}

func (tp *TicketPrioritiesController) Update(c *gin.Context) {
	id := c.Param("id")
	var priority models.TicketPriority
	if err := database.DB.First(&priority, id).Error; err != nil {
		utils.Error(c, "Priority not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	var req map[string]interface{}
	c.ShouldBindJSON(&req)
	delete(req, "id")
	database.DB.Model(&priority).Updates(req)
	utils.Success(c, priority, "Priority updated", http.StatusOK)
}

func (tp *TicketPrioritiesController) Delete(c *gin.Context) {
	id := c.Param("id")
	database.DB.Delete(&models.TicketPriority{}, id)
	utils.Success(c, nil, "Priority deleted", http.StatusOK)
}

// TicketMessagesController
type AdminTicketMessagesController struct{}

func (tm *AdminTicketMessagesController) Index(c *gin.Context) {
	uuid := c.Param("uuid")
	var ticket models.Ticket
	if err := database.DB.Where("uuid = ?", uuid).First(&ticket).Error; err != nil {
		utils.Error(c, "Ticket not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	var msgs []models.TicketMessage
	database.DB.Where("ticket_id = ?", ticket.ID).Find(&msgs)
	utils.Success(c, msgs, "Messages retrieved", http.StatusOK)
}

func (tm *AdminTicketMessagesController) Show(c *gin.Context) {
	id := c.Param("id")
	var msg models.TicketMessage
	if err := database.DB.First(&msg, id).Error; err != nil {
		utils.Error(c, "Message not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	utils.Success(c, msg, "Message retrieved", http.StatusOK)
}

func (tm *AdminTicketMessagesController) Create(c *gin.Context) {
	uuid := c.Param("uuid")
	var ticket models.Ticket
	if err := database.DB.Where("uuid = ?", uuid).First(&ticket).Error; err != nil {
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
	ctxUser, _ := c.Get("user")
	user := ctxUser.(models.User)
	msg := models.TicketMessage{TicketID: int(ticket.ID), UserUUID: &user.UUID, Message: req.Message, IsStaff: "true"}
	database.DB.Create(&msg)
	utils.Success(c, msg, "Message created", http.StatusCreated)
}

func (tm *AdminTicketMessagesController) Update(c *gin.Context) {
	id := c.Param("id")
	var msg models.TicketMessage
	if err := database.DB.First(&msg, id).Error; err != nil {
		utils.Error(c, "Message not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	var req map[string]interface{}
	c.ShouldBindJSON(&req)
	delete(req, "id")
	database.DB.Model(&msg).Updates(req)
	utils.Success(c, msg, "Message updated", http.StatusOK)
}

func (tm *AdminTicketMessagesController) Delete(c *gin.Context) {
	id := c.Param("id")
	database.DB.Delete(&models.TicketMessage{}, id)
	utils.Success(c, nil, "Message deleted", http.StatusOK)
}
