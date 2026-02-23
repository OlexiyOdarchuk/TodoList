package routes

import (
	"net/http"
	"strings"
	"todolist/internal/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(manager *utils.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)

		claims, err := manager.Validate(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "details": err.Error()})
			c.Abort()
			return
		}
		c.Set("user", claims.UserID)
		c.Next()
	}
}
