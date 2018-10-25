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
	frontEndURL, ok := os.LookupEnv("FRONTEND_URL")
	if ok {
		config.AllowOrigins = append(config.AllowOrigins, frontEndURL)
	}
	r.Use(cors.New(config))
	api := r.Group("/api")
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": true})
		})
		api.POST("/login", login)
	}

	r.GET("/.well-known/jwks.json", jwks)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(fmt.Sprintf(":%s", port))
}
