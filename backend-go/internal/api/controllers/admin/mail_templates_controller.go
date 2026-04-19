package admin

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mythicalsystems/nemesis-backend/internal/database"
	"github.com/mythicalsystems/nemesis-backend/internal/models"
	"github.com/mythicalsystems/nemesis-backend/pkg/utils"
)

type MailTemplatesController struct{}

func (m *MailTemplatesController) Index(c *gin.Context) {
	p := utils.GetPagination(c)
	search := strings.TrimSpace(c.Query("search"))

	q := database.DB.Model(&models.MailTemplate{})
	if search != "" {
		like := "%" + search + "%"
		q = q.Where("name LIKE ? OR subject LIKE ?", like, like)
	}

	var total int64
	q.Count(&total)

	var templates []models.MailTemplate
	q.Order("created_at DESC").Offset(p.Offset).Limit(p.PerPage).Find(&templates)
	if templates == nil {
		templates = []models.MailTemplate{}
	}

	utils.Success(c, gin.H{
		"templates":  templates,
		"pagination": utils.BuildPagination(p, total),
	}, "Mail templates retrieved", http.StatusOK)
}

func (m *MailTemplatesController) Show(c *gin.Context) {
	id := c.Param("id")
	var tpl models.MailTemplate
	if err := database.DB.First(&tpl, id).Error; err != nil {
		utils.Error(c, "Template not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	utils.Success(c, gin.H{"template": tpl}, "Template retrieved", http.StatusOK)
}

func (m *MailTemplatesController) Create(c *gin.Context) {
	var req struct {
		Name    string `json:"name" binding:"required"`
		Subject string `json:"subject" binding:"required"`
		Body    string `json:"body" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	tpl := models.MailTemplate{Name: req.Name, Subject: req.Subject, Body: req.Body}
	if err := database.DB.Create(&tpl).Error; err != nil {
		utils.Error(c, "A template with that name already exists", "DUPLICATE", http.StatusConflict, nil)
		return
	}
	utils.Success(c, gin.H{"template": tpl, "id": tpl.ID}, "Template created", http.StatusCreated)
}

func (m *MailTemplatesController) Update(c *gin.Context) {
	id := c.Param("id")
	var tpl models.MailTemplate
	if err := database.DB.First(&tpl, id).Error; err != nil {
		utils.Error(c, "Template not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	delete(req, "id")
	delete(req, "created_at")
	delete(req, "updated_at")
	database.DB.Model(&tpl).Updates(req)
	utils.Success(c, gin.H{"template": tpl}, "Template updated", http.StatusOK)
}

func (m *MailTemplatesController) Delete(c *gin.Context) {
	id := c.Param("id")
	database.DB.Delete(&models.MailTemplate{}, id)
	utils.Success(c, nil, "Template deleted", http.StatusOK)
}

func (m *MailTemplatesController) TestEmail(c *gin.Context) {
	var req struct {
		Email   string `json:"email" binding:"required,email"`
		Subject string `json:"subject" binding:"required"`
		Body    string `json:"body" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	database.DB.Create(&models.MailQueue{
		Subject: req.Subject,
		Body:    req.Body,
		Status:  "pending",
	})
	utils.Success(c, gin.H{"queued": 1, "email": req.Email}, "Test email queued", http.StatusOK)
}

func (m *MailTemplatesController) MassEmail(c *gin.Context) {
	var req struct {
		Subject string `json:"subject" binding:"required"`
		Body    string `json:"body" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}

	var users []models.User
	database.DB.Where("banned = ?", "false").Find(&users)

	queued := 0
	for _, u := range users {
		if err := database.DB.Create(&models.MailQueue{
			UserUUID: u.UUID,
			Subject:  req.Subject,
			Body:     req.Body,
			Status:   "pending",
		}).Error; err == nil {
			queued++
		}
	}
	utils.Success(c, gin.H{"queued_count": queued}, "Mass email queued", http.StatusOK)
}
