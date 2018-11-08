package web

import (
	"net/http"

	"github.com/bottleneckco/presence-backend/model"
	"github.com/gin-gonic/gin"
)

func me(c *gin.Context) {
	user := c.MustGet("user").(model.User)
	c.JSON(http.StatusOK, gin.H{"status": true, "data": user})
}
