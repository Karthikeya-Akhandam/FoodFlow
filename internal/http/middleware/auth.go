package middleware

import (
	"strings"

	"foodflow/internal/core"
	"foodflow/internal/lib"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	jwtManager *lib.JWTManager
}

func NewAuthMiddleware(jwtManager *lib.JWTManager) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager: jwtManager,
	}
}

func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			lib.NewResponder().Unauthorized(c, "Authorization header required")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			lib.NewResponder().Unauthorized(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := m.jwtManager.ValidateToken(token)
		if err != nil {
			lib.NewResponder().Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		// Set user context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

func (m *AuthMiddleware) RequireRole(roles ...core.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			lib.NewResponder().Unauthorized(c, "User role not found in context")
			c.Abort()
			return
		}

		role, ok := userRole.(core.UserRole)
		if !ok {
			lib.NewResponder().Unauthorized(c, "Invalid user role type")
			c.Abort()
			return
		}

		hasRole := false
		for _, allowedRole := range roles {
			if role == allowedRole {
				hasRole = true
				break
			}
		}

		if !hasRole {
			lib.NewResponder().Forbidden(c, "Insufficient permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

func (m *AuthMiddleware) RequireAdmin() gin.HandlerFunc {
	return m.RequireRole(core.UserRoleAdmin)
}

func (m *AuthMiddleware) RequireOrg() gin.HandlerFunc {
	return m.RequireRole(core.UserRoleOrg)
}

func (m *AuthMiddleware) RequireCollab() gin.HandlerFunc {
	return m.RequireRole(core.UserRoleCollab)
}

func (m *AuthMiddleware) RequireOrgOrCollab() gin.HandlerFunc {
	return m.RequireRole(core.UserRoleOrg, core.UserRoleCollab)
}

// GetUserID extracts user ID from context
func GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}
	id, ok := userID.(string)
	return id, ok
}

// GetUserRole extracts user role from context
func GetUserRole(c *gin.Context) (core.UserRole, bool) {
	userRole, exists := c.Get("user_role")
	if !exists {
		return "", false
	}
	role, ok := userRole.(core.UserRole)
	return role, ok
}
