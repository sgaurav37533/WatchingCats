package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gaurav/watchingcat/internal/config"
)

// CORS returns a gin middleware for handling CORS
func CORS(cfg config.CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Set CORS headers
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			// Check if origin is allowed
			allowed := false
			for _, allowedOrigin := range cfg.AllowedOrigins {
				if allowedOrigin == "*" || allowedOrigin == origin {
					allowed = true
					break
				}
			}

			if allowed {
				c.Header("Access-Control-Allow-Origin", origin)
			}
		}

		// Set other CORS headers
		c.Header("Access-Control-Allow-Methods", join(cfg.AllowedMethods, ", "))
		c.Header("Access-Control-Allow-Headers", join(cfg.AllowedHeaders, ", "))
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Max-Age", "86400")

		// Handle preflight request
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func join(items []string, sep string) string {
	if len(items) == 0 {
		return ""
	}
	result := items[0]
	for _, item := range items[1:] {
		result += sep + item
	}
	return result
}

