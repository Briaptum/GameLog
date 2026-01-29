package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Auth0Claims represents the JWT claims from Auth0
type Auth0Claims struct {
	Sub           string                 `json:"sub"`
	Email         string                 `json:"email"`
	AppMetadata   map[string]interface{} `json:"https://your-namespace.com/app_metadata"`
	jwt.RegisteredClaims
}

// Auth0Middleware handles Auth0 JWT validation and user authorization
func Auth0Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		tokenString := extractToken(c)
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		// Parse and validate the token
		token, err := jwt.ParseWithClaims(tokenString, &Auth0Claims{}, func(token *jwt.Token) (interface{}, error) {
			// Verify the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// Use the same secret key as in the login controller
			return []byte("your-secret-key"), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "details": err.Error()})
			c.Abort()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(*Auth0Claims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Check authorization based on app metadata
		if !isAuthorized(claims) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		// Store user info in context
		c.Set("user", claims)
		c.Next()
	}
}

// extractToken extracts the JWT token from the Authorization header
func extractToken(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return ""
	}

	// Check if it's a Bearer token
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

// isAuthorized checks if the user has permission to access the admin panel
func isAuthorized(claims *Auth0Claims) bool {
	// Check if user has admin role
	if role, exists := claims.AppMetadata["role"]; exists {
		if roleStr, ok := role.(string); ok && roleStr == "admin" {
			return true
		}
	}

	// Check if user has website-dummy-com in their sites list
	if sites, exists := claims.AppMetadata["sites"]; exists {
		if sitesList, ok := sites.([]interface{}); ok {
			for _, site := range sitesList {
				if siteStr, ok := site.(string); ok && siteStr == "website-dummy-com" {
					return true
				}
			}
		}
	}

	return false
}

// Optional: Middleware for public routes that don't require authentication
func OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractToken(c)
		if tokenString == "" {
			c.Next()
			return
		}

		// Try to parse token, but don't fail if invalid
		token, err := jwt.ParseWithClaims(tokenString, &Auth0Claims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("AUTH0_CLIENT_SECRET")), nil
		})

		if err == nil && token.Valid {
			if claims, ok := token.Claims.(*Auth0Claims); ok {
				c.Set("user", claims)
			}
		}

		c.Next()
	}
}
