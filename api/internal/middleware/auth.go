package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// APIKeyAuth valida o header Authorization: Bearer <api_key>.
func APIKeyAuth(validKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing Authorization header"})
			c.Abort()
			return
		}

		parts := strings.SplitN(auth, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid Authorization format, use Bearer <token>"})
			c.Abort()
			return
		}

		key := strings.TrimSpace(parts[1])
		if key != validKey {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid API key"})
			c.Abort()
			return
		}

		c.Next()
	}
}
