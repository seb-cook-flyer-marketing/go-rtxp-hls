package main

import (
	"github.com/gin-gonic/gin"
	"github.com/seb-cook-flyer-marketing/go-rtxp-hls/config"
	"github.com/seb-cook-flyer-marketing/go-rtxp-hls/stream"
)

func main() {
	// Load configuration
	cfg := config.Config

	// Initialize Gin router
	r := gin.Default()

	r.Static("/static", "./public")
	// Example route
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to the API!",
			"port":    cfg.Port,
		})
	})

	stream.RegisterStreamRoutes(r)

	// Start the server
	r.Run(":" + cfg.Port) // Listen and serve on the specified port
}
