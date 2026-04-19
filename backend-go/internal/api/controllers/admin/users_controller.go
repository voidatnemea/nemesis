package admin

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mythicalsystems/nemesis-backend/internal/database"
	"github.com/mythicalsystems/nemesis-backend/internal/models"
	"github.com/mythicalsystems/nemesis-backend/pkg/utils"
)

type UsersController struct{}

func (u *UsersController) Index(c *gin.Context) {
	p := utils.GetPagination(c)
	search := c.Query("search")

	var users []models.User
	var total int64
	q := database.DB.Model(&models.User{})
	if search != "" {
		like := "%" + search + "%"
		q = q.Where("username LIKE ? OR email LIKE ? OR first_name LIKE ? OR last_name LIKE ?", like, like, like, like)
	}
	if role := c.Query("role"); role != "" {
		q = q.Where("role_id = ?", role)
	}
	if banned := c.Query("banned"); banned != "" {
		q = q.Where("banned = ?", banned)
	}
	q.Count(&total)
	q.Offset(p.Offset).Limit(p.PerPage).Find(&users)

	var roles []models.Role
	database.DB.Find(&roles)
	rolesMap := make(map[uint]models.Role, len(roles))
	for _, r := range roles {
		rolesMap[r.ID] = r
	}

	utils.Success(c, gin.H{
		"users":      users,
		"roles":      rolesMap,
		"pagination": utils.BuildPagination(p, total),
	}, "Users retrieved", http.StatusOK)
}

func (u *UsersController) Show(c *gin.Context) {
	uuid := c.Param("uuid")
	var user models.User
	if err := database.DB.Where("uuid = ?", uuid).First(&user).Error; err != nil {
		utils.Error(c, "User not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	var role models.Role
	database.DB.First(&role, user.RoleID)
	utils.Success(c, gin.H{"user": user, "role": role}, "User retrieved", http.StatusOK)
}

func (u *UsersController) ShowByExternalID(c *gin.Context) {
	externalID := c.Param("externalId")
	var user models.User
	if err := database.DB.Where("external_id = ?", externalID).First(&user).Error; err != nil {
		utils.Error(c, "User not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	utils.Success(c, user, "User retrieved", http.StatusOK)
}

func (u *UsersController) Create(c *gin.Context) {
	var req struct {
		Username  string `json:"username" binding:"required,min=3,max=50"`
		Email     string `json:"email" binding:"required,email"`
		Password  string `json:"password" binding:"required,min=8"`
		FirstName string `json:"first_name" binding:"required"`
		LastName  string `json:"last_name" binding:"required"`
		RoleID    int    `json:"role_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	var count int64
	database.DB.Model(&models.User{}).Where("username = ? OR email = ?", req.Username, req.Email).Count(&count)
	if count > 0 {
		utils.Error(c, "Username or email already exists", "CONFLICT", http.StatusConflict, nil)
		return
	}
	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.Error(c, "Failed to hash password", "INTERNAL", http.StatusInternalServerError, nil)
		return
	}
	roleID := req.RoleID
	if roleID == 0 {
		roleID = 2
	}
	user := models.User{
		UUID:      utils.GenerateUUID(),
		Username:  req.Username,
		Email:     req.Email,
		Password:  hashed,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		RoleID:    roleID,
	}
	if err := database.DB.Create(&user).Error; err != nil {
		utils.Error(c, "Failed to create user", "DATABASE_ERROR", http.StatusInternalServerError, nil)
		return
	}
	utils.Success(c, user, "User created", http.StatusCreated)
}

func (u *UsersController) Update(c *gin.Context) {
	uuid := c.Param("uuid")
	var user models.User
	if err := database.DB.Where("uuid = ?", uuid).First(&user).Error; err != nil {
		utils.Error(c, "User not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	if pw, ok := req["password"].(string); ok && pw != "" {
		hashed, _ := utils.HashPassword(pw)
		req["password"] = hashed
	}
	delete(req, "uuid")
	delete(req, "id")
	database.DB.Model(&user).Updates(req)
	utils.Success(c, user, "User updated", http.StatusOK)
}

func (u *UsersController) Delete(c *gin.Context) {
	uuid := c.Param("uuid")
	if err := database.DB.Where("uuid = ?", uuid).Delete(&models.User{}).Error; err != nil {
		utils.Error(c, "Failed to delete user", "DATABASE_ERROR", http.StatusInternalServerError, nil)
		return
	}
	utils.Success(c, nil, "User deleted", http.StatusOK)
}

func (u *UsersController) OwnedServers(c *gin.Context) {
	uuid := c.Param("uuid")
	var user models.User
	if err := database.DB.Where("uuid = ?", uuid).First(&user).Error; err != nil {
		utils.Error(c, "User not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	var servers []models.Server
	database.DB.Where("owner_id = ?", user.ID).Find(&servers)
	utils.Success(c, gin.H{"servers": servers}, "Servers retrieved", http.StatusOK)
}

func (u *UsersController) Ban(c *gin.Context) {
	uuid := c.Param("uuid")
	database.DB.Model(&models.User{}).Where("uuid = ?", uuid).Update("banned", "true")
	utils.Success(c, nil, "User banned", http.StatusOK)
}

func (u *UsersController) Unban(c *gin.Context) {
	uuid := c.Param("uuid")
	database.DB.Model(&models.User{}).Where("uuid = ?", uuid).Update("banned", "false")
	utils.Success(c, nil, "User unbanned", http.StatusOK)
}

func (u *UsersController) CreateSsoToken(c *gin.Context) {
	uuid := c.Param("uuid")
	var user models.User
	if err := database.DB.Where("uuid = ?", uuid).First(&user).Error; err != nil {
		utils.Error(c, "User not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	token := utils.GenerateRandomString(32)
	sso := models.SsoToken{
		UserUUID:  user.UUID,
		Token:     token,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
	database.DB.Create(&sso)
	utils.Success(c, gin.H{"token": token}, "SSO token created", http.StatusCreated)
}
