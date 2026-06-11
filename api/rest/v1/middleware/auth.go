package middleware

import (
	"construction_transport_server/api/rest/v1/delivery"
	"construction_transport_server/pkg/utils"
	"github.com/gin-gonic/gin"
	"strings"
)

func AuthMiddleware(jwtManager *utils.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			delivery.SendError(c, 401, "missing authorization header")
			c.Abort()
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			delivery.SendError(c, 401, "invalid authorization format")
			c.Abort()
			return
		}
		token := parts[1]
		claims, err := jwtManager.ValidateAccessToken(token)
		if err != nil {
			delivery.SendError(c, 401, "invalid or expired token")
			c.Abort()
			return
		}
		c.Set("auth_id", claims.AuthId)
		c.Set("role", claims.Role)
		c.Set("email", claims.Email)
		c.Next()
	}
}

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		for _, r := range allowedRoles {
			if r == role {
				c.Next()
				return
			}
		}
		delivery.SendError(c, 403, "insufficient permissions")
		c.Abort()
	}
}
