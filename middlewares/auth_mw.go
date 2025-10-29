package middlewares

import (
	"cineverse/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// func AuthMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		auth := c.GetHeader("Authorization")
// 		token := ""
// 		if auth == "" {
// 			// For HTML requests, check query param
// 			token = c.Query("token")
// 			if token == "" {
// 				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
// 				return
// 			}
// 		} else {
// 			parts := strings.Split(auth, " ")
// 			if len(parts) != 2 || parts[0] != "Bearer" {
// 				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
// 				return
// 			}
// 			token = parts[1]
// 		}
// 		claims, err := utils.ParseToken(token)
// 		if err != nil {
// 			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token", "details": err.Error()})
// 			return
// 		}
// 		// attach to context
// 		c.Set("user_id", claims.UserID)
// 		c.Set("is_admin", claims.IsAdmin)
// 		c.Next()
// 	}
// }

// func AdminMiddleware() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		v, exists := c.Get("is_admin")
// 		if !exists {
// 			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "no admin info"})
// 			return
// 		}
// 		isAdmin, ok := v.(bool)
// 		if !ok || !isAdmin {
// 			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin only"})
// 			return
// 		}
// 		c.Next()
// 	}
// }

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role") // role should be set inside JWT claims
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied: Admins only"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		userID, role, err := utils.ValidateJWT(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Set userID and role in context (so AdminMiddleware can use it)
		c.Set("userID", userID)
		c.Set("role", role)

		c.Next()
	}
}
