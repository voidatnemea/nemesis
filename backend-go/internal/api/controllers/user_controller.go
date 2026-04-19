package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mythicalsystems/nemesis-backend/internal/models"
	"github.com/mythicalsystems/nemesis-backend/pkg/utils"
)

type UserController struct{}

func (u *UserController) Profile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		utils.Error(c, "User not found in context", "USER_NOT_FOUND", http.StatusInternalServerError, nil)
		return
	}

	utils.Success(c, user.(models.User), "User profile retrieved successfully", http.StatusOK)
}
