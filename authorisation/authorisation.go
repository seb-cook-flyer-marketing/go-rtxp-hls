package authorisation

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/seb-cook-flyer-marketing/go-rtxp-hls/config"
)

// AuthenticationMiddleware is a middleware that checks for a valid Bearer token.
func AuthenticationMiddleware(next gin.HandlerFunc) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// No authorization on development
		if config.Config.Secret == "" {
			next(c)
			return
		}

		authorization := c.Request.Header.Get("Authorization")
		if authorization == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Expected format: "Bearer <token>"
		parts := strings.SplitN(authorization, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		token := parts[1]
		if token != config.Config.Secret {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Proceed to the next handler
		next(c)
	})
}
