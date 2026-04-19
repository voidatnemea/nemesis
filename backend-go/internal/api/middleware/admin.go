package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mythicalsystems/nemesis-backend/internal/database"
	"github.com/mythicalsystems/nemesis-backend/internal/models"
	"github.com/mythicalsystems/nemesis-backend/pkg/utils"
)

// AdminMiddleware verifies that the authenticated user holds a permission.
// Call after AuthMiddleware. Pass an empty string to only require admin role (role_id == 1).
func AdminMiddleware(requiredPermission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctxUser, exists := c.Get("user")
		if !exists {
			utils.Error(c, "Unauthorized", "UNAUTHORIZED", http.StatusUnauthorized, nil)
			c.Abort()
			return
		}
		user := ctxUser.(models.User)

		// role_id 1 is the root/super-admin role
		if user.RoleID == 1 {
			c.Next()
			return
		}

		if requiredPermission == "" {
			utils.Error(c, "You do not have permission to access this resource", "FORBIDDEN", http.StatusForbidden, nil)
			c.Abort()
			return
		}

		// Check specific permission
		var perm models.Permission
		err := database.DB.
			Joins("JOIN featherpanel_roles ON featherpanel_roles.id = featherpanel_permissions.role_id").
			Where("featherpanel_permissions.role_id = ? AND featherpanel_permissions.permission = ?", user.RoleID, requiredPermission).
			First(&perm).Error
		if err != nil {
			utils.Error(c, "You do not have permission to access this resource", "FORBIDDEN", http.StatusForbidden, nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
