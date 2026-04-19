package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mythicalsystems/nemesis-backend/internal/database"
	"github.com/mythicalsystems/nemesis-backend/internal/models"
	"github.com/mythicalsystems/nemesis-backend/pkg/utils"
)

type SessionController struct{}

func (s *SessionController) Get(c *gin.Context) {
	ctxUser, _ := c.Get("user")
	user := ctxUser.(models.User)
	var role models.Role
	database.DB.First(&role, user.RoleID)
	var perms []models.Permission
	database.DB.Where("role_id = ?", user.RoleID).Find(&perms)
	var notifs []models.Notification
	database.DB.Order("created_at DESC").Limit(5).Find(&notifs)
	utils.Success(c, gin.H{
		"user":          user,
		"role":          role,
		"permissions":   perms,
		"notifications": notifs,
	}, "Session retrieved", http.StatusOK)
}

func (s *SessionController) Update(c *gin.Context) {
	ctxUser, _ := c.Get("user")
	user := ctxUser.(models.User)
	var req struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	updates := map[string]interface{}{}
	if req.FirstName != "" {
		updates["first_name"] = req.FirstName
	}
	if req.LastName != "" {
		updates["last_name"] = req.LastName
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	database.DB.Model(&user).Updates(updates)
	utils.Success(c, user, "Session updated", http.StatusOK)
}
