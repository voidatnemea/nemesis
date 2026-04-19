package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mythicalsystems/nemesis-backend/internal/database"
	"github.com/mythicalsystems/nemesis-backend/internal/models"
	"github.com/mythicalsystems/nemesis-backend/pkg/utils"
)

type UserActivitiesController struct{}

func (a *UserActivitiesController) Index(c *gin.Context) {
	ctxUser, _ := c.Get("user")
	user := ctxUser.(models.User)

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if limit < 1 || limit > 100 {
		limit = 20
	}

	var activities []models.Activity
	database.DB.Where("user_uuid = ?", user.UUID).Order("created_at DESC").Limit(limit).Find(&activities)

	utils.Success(c, gin.H{
		"activities": activities,
	}, "Activities retrieved", http.StatusOK)
}
