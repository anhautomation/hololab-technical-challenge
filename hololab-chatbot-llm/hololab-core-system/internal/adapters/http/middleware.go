package httpadapter

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		if origin != "" && !isAllowedOrigin(origin) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Not Allowed"})
			return
		}

		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Vary", "Origin")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Max-Age", "86400")
		}

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func isAllowedOrigin(origin string) bool {
	if origin == "" {
		return false
	}

	allowedExact := []string{
		"https://hololab-technical-challenge.vercel.app",
		"http://localhost:5173",
		"http://localhost:3000",
	}

	for _, o := range allowedExact {
		if origin == o {
			return true
		}
	}

	return false
}
