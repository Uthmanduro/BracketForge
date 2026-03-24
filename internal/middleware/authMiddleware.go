package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/uthmanduro/BracketForge/internal/auth"
	"github.com/uthmanduro/BracketForge/internal/config"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract JWT token from the request
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Authorization header missing"})
			return
		}

		// Split the token string to get the actual token (assuming "Bearer <token>")
		tokenParts := strings.Split(tokenString, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid Authorization header format"})
			return
		}
		tokenString = tokenParts[1]

		// Validate the JWT token and extract user information
		claims, err := auth.ValidateJWT(tokenString, cfg.JWTSecret)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid or expired token"})
			return
		}

		// Store user information in the context for use in handlers
		c.Set("userID", claims.UserID)
		c.Set("userRole", claims.Role)
		c.Set("organizationID", claims.OrganizationID)
		
		// Continue to the next handler
		c.Next()
	}
}

func RequireRole(roles ...string) gin.HandlerFunc {
	allowed := make(map[string]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		if !allowed[role.(string)] {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			return
		}
		c.Next()
	}
}
 
func OrgID(c *gin.Context) string {
	v, _ := c.Get("organization_id")
	s, _ := v.(string)
	return s
}