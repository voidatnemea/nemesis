package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mythicalsystems/nemesis-backend/internal/database"
	"github.com/mythicalsystems/nemesis-backend/internal/models"
	"github.com/mythicalsystems/nemesis-backend/pkg/utils"
)

type UserSshKeyController struct{}

func (s *UserSshKeyController) GetUserSshKeys(c *gin.Context) {
	ctxUser, _ := c.Get("user")
	user := ctxUser.(models.User)
	var keys []models.UserSshKey
	database.DB.Where("user_id = ?", user.ID).Find(&keys)
	utils.Success(c, keys, "SSH keys retrieved", http.StatusOK)
}

func (s *UserSshKeyController) GetUserSshKey(c *gin.Context) {
	ctxUser, _ := c.Get("user")
	user := ctxUser.(models.User)
	id := c.Param("id")
	var key models.UserSshKey
	if err := database.DB.Where("id = ? AND user_id = ?", id, user.ID).First(&key).Error; err != nil {
		utils.Error(c, "SSH key not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	utils.Success(c, key, "SSH key retrieved", http.StatusOK)
}

func (s *UserSshKeyController) CreateUserSshKey(c *gin.Context) {
	ctxUser, _ := c.Get("user")
	user := ctxUser.(models.User)
	var req struct {
		Name      string `json:"name" binding:"required"`
		PublicKey string `json:"public_key" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	key := models.UserSshKey{
		UserID:    int(user.ID),
		Name:      req.Name,
		PublicKey: req.PublicKey,
	}
	database.DB.Create(&key)
	utils.Success(c, key, "SSH key created", http.StatusCreated)
}

func (s *UserSshKeyController) UpdateUserSshKey(c *gin.Context) {
	ctxUser, _ := c.Get("user")
	user := ctxUser.(models.User)
	id := c.Param("id")
	var key models.UserSshKey
	if err := database.DB.Where("id = ? AND user_id = ?", id, user.ID).First(&key).Error; err != nil {
		utils.Error(c, "SSH key not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	var req struct {
		Name string `json:"name"`
	}
	c.ShouldBindJSON(&req)
	if req.Name != "" {
		database.DB.Model(&key).Update("name", req.Name)
	}
	utils.Success(c, key, "SSH key updated", http.StatusOK)
}

func (s *UserSshKeyController) DeleteUserSshKey(c *gin.Context) {
	ctxUser, _ := c.Get("user")
	user := ctxUser.(models.User)
	id := c.Param("id")
	database.DB.Where("id = ? AND user_id = ?", id, user.ID).Delete(&models.UserSshKey{})
	utils.Success(c, nil, "SSH key deleted", http.StatusOK)
}
