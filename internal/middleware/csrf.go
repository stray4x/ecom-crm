package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CSRFMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		csrfCookie, err := c.Cookie("csrf_token")
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "csrf token missing"})
			c.Abort()
			return
		}

		csrfHeader := c.GetHeader("X-CSRF-Token")
		if csrfHeader == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "csrf token missing"})
			c.Abort()
			return
		}

		if csrfCookie != csrfHeader {
			c.JSON(http.StatusForbidden, gin.H{"error": "csrf token invalid"})
			c.Abort()
			return
		}

		c.Next()
	}
}
