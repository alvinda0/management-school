package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"project-go/config"
	"project-go/internal/model"
	"project-go/utils"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Authorization header required")
			c.Abort()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		if tokenString == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid token format")
			c.Abort()
			return
		}

		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid token")
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("roleID", claims.RoleID)
		c.Next()
	}
}

func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleID, exists := c.Get("roleID")
		if !exists {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Role not found in token")
			c.Abort()
			return
		}

		// Check if the role has the required permission
		hasPermission := checkRolePermission(roleID.(uint), permission)
		if !hasPermission {
			utils.ErrorResponse(c, http.StatusForbidden, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

// Dynamic permission check using database
func checkRolePermission(roleID uint, permission string) bool {
	// Get role with permissions from database
	var role model.Role
	err := config.DB.Preload("Permissions").First(&role, roleID).Error
	if err != nil {
		return false
	}

	// System and Kepala Sekolah have full access to everything
	if role.Name == "system" || role.Name == "kepala_sekolah" {
		return true
	}

	// Check if role has the specific permission
	for _, perm := range role.Permissions {
		if perm.Name == permission {
			return true
		}
	}

	return false
}

// RequireRole checks if user has specific role
func RequireRole(roleName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleID, exists := c.Get("roleID")
		if !exists {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Role not found in token")
			c.Abort()
			return
		}

		// Get role from database
		var role model.Role
		err := config.DB.First(&role, roleID.(uint)).Error
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid role")
			c.Abort()
			return
		}

		// System and Kepala Sekolah can access everything
		if role.Name == "system" || role.Name == "kepala_sekolah" {
			c.Next()
			return
		}

		// Check if user has the required role
		if role.Name != roleName {
			utils.ErrorResponse(c, http.StatusForbidden, "Insufficient role privileges")
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyRole checks if user has any of the specified roles
func RequireAnyRole(roleNames ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleID, exists := c.Get("roleID")
		if !exists {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Role not found in token")
			c.Abort()
			return
		}

		// Get role from database
		var role model.Role
		err := config.DB.First(&role, roleID.(uint)).Error
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid role")
			c.Abort()
			return
		}

		// System and Kepala Sekolah can access everything
		if role.Name == "system" || role.Name == "kepala_sekolah" {
			c.Next()
			return
		}

		// Check if user has any of the required roles
		for _, roleName := range roleNames {
			if role.Name == roleName {
				c.Next()
				return
			}
		}

		utils.ErrorResponse(c, http.StatusForbidden, "Insufficient role privileges")
		c.Abort()
	}
}

// IsSystemOrKepalaSekolah checks if user is system admin or kepala sekolah
func IsSystemOrKepalaSekolah() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleID, exists := c.Get("roleID")
		if !exists {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Role not found in token")
			c.Abort()
			return
		}

		// Get role from database
		var role model.Role
		err := config.DB.First(&role, roleID.(uint)).Error
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid role")
			c.Abort()
			return
		}

		// Only system and kepala sekolah can access
		if role.Name != "system" && role.Name != "kepala_sekolah" {
			utils.ErrorResponse(c, http.StatusForbidden, "Access restricted to system admin and kepala sekolah only")
			c.Abort()
			return
		}

		c.Next()
	}
}