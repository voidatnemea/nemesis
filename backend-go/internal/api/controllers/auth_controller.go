package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mythicalsystems/nemesis-backend/internal/database"
	"github.com/mythicalsystems/nemesis-backend/internal/models"
	"github.com/mythicalsystems/nemesis-backend/pkg/utils"
)

type AuthController struct{}

type LoginRequest struct {
	UsernameOrEmail string `json:"username_or_email" binding:"required"`
	Password        string `json:"password" binding:"required"`
}

func (a *AuthController) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, "Invalid request payload", "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}

	var user models.User
	// Try email first
	if err := database.DB.Where("email = ?", req.UsernameOrEmail).First(&user).Error; err != nil {
		// Try username
		if err := database.DB.Where("username = ?", req.UsernameOrEmail).First(&user).Error; err != nil {
			utils.Error(c, "Invalid username or email address", "INVALID_USERNAME_OR_EMAIL", http.StatusUnauthorized, nil)
			return
		}
	}

	if user.Banned == "true" {
		utils.Error(c, "User is banned", "USER_BANNED", http.StatusForbidden, nil)
		return
	}

	if !utils.VerifyPassword(req.Password, user.Password) {
		utils.Error(c, "Invalid password", "INVALID_PASSWORD", http.StatusUnauthorized, nil)
		return
	}

	// 2FA check
	if user.TwoFAEnabled == "true" {
		utils.Error(c, "2FA required", "TWO_FACTOR_REQUIRED", http.StatusUnauthorized, gin.H{
			"email": user.Email,
		})
		return
	}

	// Rotate remember token on every login to reduce token replay risk
	rawToken, hashedToken, err := utils.GenerateSessionToken()
	if err != nil {
		utils.Error(c, "Failed to generate session token", "TOKEN_GENERATION_FAILED", http.StatusInternalServerError, nil)
		return
	}
	user.RememberToken = hashedToken
	if err := database.DB.Save(&user).Error; err != nil {
		utils.Error(c, "Failed to persist session token", "DATABASE_ERROR", http.StatusInternalServerError, nil)
		return
	}

	// Set cookie — httpOnly=false so the client-side JS auth guard can read it
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("remember_token", rawToken, 3600*24*30, "/", "", c.Request.TLS != nil, false)

	utils.Success(c, gin.H{
		"user": user,
	}, "User logged in successfully", http.StatusOK)
}

type RegisterRequest struct {
	Username  string `json:"username" binding:"required,min=3,max=50"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8,max=128"`
	FirstName string `json:"first_name" binding:"required,min=1,max=100"`
	LastName  string `json:"last_name" binding:"required,min=1,max=100"`
}

func (a *AuthController) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, "Invalid request payload", "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}

	// Uniqueness check
	var count int64
	database.DB.Model(&models.User{}).Where("username = ? OR email = ?", req.Username, req.Email).Count(&count)
	if count > 0 {
		utils.Error(c, "Username or email already exists", "CONFLICT", http.StatusConflict, nil)
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.Error(c, "Failed to hash password", "PASSWORD_HASH_ERROR", http.StatusInternalServerError, nil)
		return
	}
	rawToken, hashedToken, err := utils.GenerateSessionToken()
	if err != nil {
		utils.Error(c, "Failed to generate session token", "TOKEN_GENERATION_FAILED", http.StatusInternalServerError, nil)
		return
	}
	user := models.User{
		Username:      req.Username,
		Email:         req.Email,
		Password:      hashedPassword,
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		UUID:          utils.GenerateUUID(),
		RememberToken: hashedToken,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		utils.Error(c, "Failed to create user", "DATABASE_ERROR", http.StatusInternalServerError, nil)
		return
	}

	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("remember_token", rawToken, 3600*24*30, "/", "", c.Request.TLS != nil, false)

	utils.Success(c, gin.H{
		"user": user,
	}, "User registered successfully and logged in", http.StatusOK)
}

func (a *AuthController) Logout(c *gin.Context) {
	// Best-effort token revocation
	if ctxUser, exists := c.Get("user"); exists {
		if user, ok := ctxUser.(models.User); ok {
			database.DB.Model(&models.User{}).Where("id = ?", user.ID).Update("remember_token", "")
		}
	}

	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("remember_token", "", -1, "/", "", c.Request.TLS != nil, false)
	utils.Success(c, nil, "User logged out successfully", http.StatusOK)
}
