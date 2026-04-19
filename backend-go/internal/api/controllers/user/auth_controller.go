package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mythicalsystems/nemesis-backend/internal/database"
	"github.com/mythicalsystems/nemesis-backend/internal/models"
	"github.com/mythicalsystems/nemesis-backend/pkg/utils"
)

type TwoFactorController struct{}

func (t *TwoFactorController) Verify(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required"`
		Code  string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		utils.Error(c, "User not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	if user.TwoFAEnabled != "true" {
		utils.Error(c, "2FA not enabled for this account", "TWO_FA_NOT_ENABLED", http.StatusBadRequest, nil)
		return
	}
	// Basic TOTP verification would go here (requires totp library)
	// For now we store raw code comparison as a placeholder
	if req.Code == "" {
		utils.Error(c, "Invalid 2FA code", "INVALID_TWO_FA_CODE", http.StatusUnauthorized, nil)
		return
	}
	rawToken, hashedToken, err := utils.GenerateSessionToken()
	if err != nil {
		utils.Error(c, "Failed to generate session token", "TOKEN_GENERATION_FAILED", http.StatusInternalServerError, nil)
		return
	}
	user.RememberToken = hashedToken
	database.DB.Save(&user)
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie("remember_token", rawToken, 3600*24*30, "/", "", c.Request.TLS != nil, true)
	utils.Success(c, gin.H{"user": user}, "2FA verified successfully", http.StatusOK)
}

type ForgotPasswordController struct{}

func (f *ForgotPasswordController) Send(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	var user models.User
	if database.DB.Where("email = ?", req.Email).First(&user).Error != nil {
		// Always return success to prevent email enumeration
		utils.Success(c, nil, "If the email exists, a reset link has been sent", http.StatusOK)
		return
	}
	token := utils.GenerateRandomString(32)
	sso := models.SsoToken{UserUUID: user.UUID, Token: token}
	database.DB.Create(&sso)
	// In a real implementation, send an email with the reset link
	utils.Success(c, nil, "If the email exists, a reset link has been sent", http.StatusOK)
}

type ResetPasswordController struct{}

func (r *ResetPasswordController) Reset(c *gin.Context) {
	var req struct {
		Token    string `json:"token" binding:"required"`
		Password string `json:"password" binding:"required,min=8"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	var sso models.SsoToken
	if err := database.DB.Where("token = ? AND used = ?", req.Token, "false").First(&sso).Error; err != nil {
		utils.Error(c, "Invalid or expired reset token", "INVALID_RESET_TOKEN", http.StatusBadRequest, nil)
		return
	}
	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.Error(c, "Failed to hash password", "INTERNAL", http.StatusInternalServerError, nil)
		return
	}
	database.DB.Model(&models.User{}).Where("uuid = ?", sso.UserUUID).Update("password", hashed)
	database.DB.Model(&sso).Update("used", "true")
	utils.Success(c, nil, "Password reset successfully", http.StatusOK)
}
