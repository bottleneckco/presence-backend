package web

import (
	"net/http"

	"github.com/bottleneckco/statuses-backend/db"
	"github.com/bottleneckco/statuses-backend/model"
	"github.com/gin-gonic/gin"
)

func statusLatest(c *gin.Context) {
	user := c.MustGet("user").(model.User)
	var status model.Status
	err := db.DB.Where("user_id = ?", user.ID).Where("end_time IS NULL").Or("end_time >= NOW()").Order("start_time DESC").First(&status).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": "internal error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": true, "data": status})
}

func statusCreate(c *gin.Context) {
	var payload model.Status
	err := c.BindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": false, "message": "bad payload"})
		return
	}
	user := c.MustGet("user").(model.User)
	var dbModel = model.Status{
		Title:     payload.Title,
		Notes:     payload.Notes,
		StartTime: payload.StartTime,
		EndTime:   payload.EndTime,
		User:      user,
	}
	err = db.DB.Create(&dbModel).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": "internal error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": true, "data": dbModel})
}
