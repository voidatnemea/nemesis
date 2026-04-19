package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mythicalsystems/nemesis-backend/internal/database"
	"github.com/mythicalsystems/nemesis-backend/internal/models"
	"github.com/mythicalsystems/nemesis-backend/pkg/utils"
)

type RolesController struct{}

func (r *RolesController) Index(c *gin.Context) {
	var roles []models.Role
	database.DB.Find(&roles)
	utils.Success(c, roles, "Roles retrieved", http.StatusOK)
}

func (r *RolesController) Show(c *gin.Context) {
	id := c.Param("id")
	var role models.Role
	if err := database.DB.First(&role, id).Error; err != nil {
		utils.Error(c, "Role not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	var perms []models.Permission
	database.DB.Where("role_id = ?", role.ID).Find(&perms)
	utils.Success(c, gin.H{"role": role, "permissions": perms}, "Role retrieved", http.StatusOK)
}

func (r *RolesController) Create(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		DisplayName string `json:"display_name" binding:"required"`
		Color       string `json:"color"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	if req.Color == "" {
		req.Color = "#6366f1"
	}
	role := models.Role{Name: req.Name, DisplayName: req.DisplayName, Color: req.Color}
	if err := database.DB.Create(&role).Error; err != nil {
		utils.Error(c, "Failed to create role", "DATABASE_ERROR", http.StatusInternalServerError, nil)
		return
	}
	utils.Success(c, role, "Role created", http.StatusCreated)
}

func (r *RolesController) Update(c *gin.Context) {
	id := c.Param("id")
	var role models.Role
	if err := database.DB.First(&role, id).Error; err != nil {
		utils.Error(c, "Role not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	var req map[string]interface{}
	c.ShouldBindJSON(&req)
	delete(req, "id")
	database.DB.Model(&role).Updates(req)
	utils.Success(c, role, "Role updated", http.StatusOK)
}

func (r *RolesController) Delete(c *gin.Context) {
	id := c.Param("id")
	database.DB.Delete(&models.Role{}, id)
	utils.Success(c, nil, "Role deleted", http.StatusOK)
}

type PermissionsController struct{}

func (p *PermissionsController) Index(c *gin.Context) {
	roleID := c.Query("role_id")
	var perms []models.Permission
	q := database.DB.Model(&models.Permission{})
	if roleID != "" {
		q = q.Where("role_id = ?", roleID)
	}
	q.Find(&perms)
	utils.Success(c, perms, "Permissions retrieved", http.StatusOK)
}

func (p *PermissionsController) Show(c *gin.Context) {
	id := c.Param("id")
	var perm models.Permission
	if err := database.DB.First(&perm, id).Error; err != nil {
		utils.Error(c, "Permission not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	utils.Success(c, perm, "Permission retrieved", http.StatusOK)
}

func (p *PermissionsController) Create(c *gin.Context) {
	var req struct {
		RoleID     int    `json:"role_id" binding:"required"`
		Permission string `json:"permission" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, err.Error(), "INVALID_REQUEST", http.StatusBadRequest, nil)
		return
	}
	perm := models.Permission{RoleID: req.RoleID, Permission: req.Permission}
	database.DB.Create(&perm)
	utils.Success(c, perm, "Permission created", http.StatusCreated)
}

func (p *PermissionsController) Update(c *gin.Context) {
	id := c.Param("id")
	var perm models.Permission
	if err := database.DB.First(&perm, id).Error; err != nil {
		utils.Error(c, "Permission not found", "NOT_FOUND", http.StatusNotFound, nil)
		return
	}
	var req map[string]interface{}
	c.ShouldBindJSON(&req)
	delete(req, "id")
	database.DB.Model(&perm).Updates(req)
	utils.Success(c, perm, "Permission updated", http.StatusOK)
}

func (p *PermissionsController) Delete(c *gin.Context) {
	id := c.Param("id")
	database.DB.Delete(&models.Permission{}, id)
	utils.Success(c, nil, "Permission deleted", http.StatusOK)
}
