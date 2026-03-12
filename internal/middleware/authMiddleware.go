package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
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
		token, userID, err := auth.ValidateJWT(tokenString, cfg.JWTSecret)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid or expired token"})
			return
		}

		// Store user information in the context for use in handlers
		c.Set("userID", userID)
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if role, ok := claims["role"].(string); ok {
				c.Set("userRole", role)
			}
		}
		// Continue to the next handler
		
		c.Next()
	}
}

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("userRole")
		if !exists {
			c.AbortWithStatusJSON(403, gin.H{"error": "User role not found"})
			return
		}

		roleStr, ok := userRole.(string)
		if !ok {
			c.AbortWithStatusJSON(403, gin.H{"error": "Invalid user role"})
			return
		}

		for _, allowedRole := range allowedRoles {
			if roleStr == allowedRole {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(403, gin.H{"error": "Forbidden: insufficient permissions"})
	}
}