package web

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// StartServer start the web server
func StartServer() {
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:1234"}
	config.AddAllowHeaders("Authorization")
	frontEndURL, ok := os.LookupEnv("FRONTEND_URL")
	if ok {
		config.AllowOrigins = append(config.AllowOrigins, frontEndURL)
	}
	r.Use(cors.New(config))
	r.POST("/oauth/token", oauth)
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": true})
	})

	api := r.Group("/api")
	{
		api.Use(authMiddleware())
		api.GET("/users/me", me)
		api.GET("/status/latest", statusLatest)
		api.POST("/status", statusCreate)
	}

	r.GET("/.well-known/jwks.json", jwks)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(fmt.Sprintf(":%s", port))
}
