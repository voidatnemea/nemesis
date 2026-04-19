package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mythicalsystems/nemesis-backend/internal/database"
	"github.com/mythicalsystems/nemesis-backend/internal/models"
	"github.com/mythicalsystems/nemesis-backend/pkg/utils"
)

type UserMailsController struct{}

func (m *UserMailsController) Index(c *gin.Context) {
	ctxUser, _ := c.Get("user")
	user := ctxUser.(models.User)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	offset := (page - 1) * limit

	var mails []models.MailQueue
	var total int64
	database.DB.Model(&models.MailQueue{}).Where("user_uuid = ?", user.UUID).Count(&total)
	database.DB.Where("user_uuid = ?", user.UUID).Order("created_at DESC").Limit(limit).Offset(offset).Find(&mails)

	utils.Success(c, gin.H{
		"data":        mails,
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": (int(total) + limit - 1) / limit,
	}, "Mails retrieved", http.StatusOK)
}