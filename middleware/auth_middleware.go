package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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

		// Here you would check if the role has the required permission
		// For now, we'll implement a basic check
		hasPermission := checkRolePermission(roleID.(uint), permission)
		if !hasPermission {
			utils.ErrorResponse(c, http.StatusForbidden, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

// Basic permission check - in production, this should query the database
func checkRolePermission(roleID uint, permission string) bool {
	// Role 1 (admin) has all permissions
	if roleID == 1 {
		return true
	}
	
	// Role 2 (teacher) has read permissions
	if roleID == 2 && (permission == "user.read" || permission == "role.read") {
		return true
	}
	
	// Role 3 (student) has basic read permission
	if roleID == 3 && permission == "user.read" {
		return true
	}
	
	return false
}