package web

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bottleneckco/presence-backend/db"
	"github.com/bottleneckco/presence-backend/model"
	"github.com/gin-gonic/gin"
	hashids "github.com/speps/go-hashids"
)

func groupList(c *gin.Context) {
	user := c.MustGet("user").(model.User)
	var groups []model.Group
	err := db.DB.Preload("Author").Model(&user).Association("Groups").Find(&groups).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": "internal error"})
		log.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": true, "data": groups})
}

type modelGroupListStatuses struct {
	Name     string         `json:"group_name"`
	Statuses []model.Status `json:"statuses"`
}

func groupListStatuses(c *gin.Context) {
	user := c.MustGet("user").(model.User)
	var groups []model.Group
	err := db.DB.Preload("Users").Preload("Author").Model(&user).Association("Groups").Find(&groups).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": "internal error"})
		log.Println(err)
		return
	}
	response := make([]modelGroupListStatuses, len(groups))
	for i, group := range groups {
		resp := modelGroupListStatuses{
			Name: group.Name,
		}
		for _, user := range group.Users {
			status := model.Status{
				User: user,
			}
			err := db.DB.Where("user_id = ?", user.ID).Where("end_time IS NULL").Or("end_time >= NOW()").Or("NOW() BETWEEN start_time AND end_time").Order("start_time DESC").First(&status).Error
			if err != nil {
				log.Println(err)
			} else {
				resp.Statuses = append(resp.Statuses, status)
			}
		}
		response[i] = resp
	}
	c.JSON(http.StatusOK, gin.H{"status": true, "data": response})
}

func groupCreate(c *gin.Context) {
	var payload model.Group
	err := c.BindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": false, "message": "bad payload"})
		return
	}
	user := c.MustGet("user").(model.User)

	hd := hashids.NewData()
	hd.Salt = fmt.Sprintf("%d %s", time.Now().UnixNano(), user.Email)
	hd.MinLength = 5

	h, _ := hashids.NewWithData(hd)

	var dbModel = model.Group{
		Name:   payload.Name,
		Code:   "",
		Author: user,
	}
	err = db.DB.Create(&dbModel).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": "internal error"})
		log.Println(err)
		return
	}

	err = db.DB.Model(&dbModel).Association("Users").Append(&user).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": "internal error"})
		log.Println(err)
		return
	}

	code, _ := h.Encode([]int{int(dbModel.ID)})
	err = db.DB.Model(&dbModel).Update("code", code).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": "internal error"})
		log.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": true, "data": dbModel})
}

type payloadGroupJoin struct {
	Code string `json:"code"`
}

func groupJoin(c *gin.Context) {
	var payload payloadGroupJoin
	err := c.BindJSON(&payload)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"status": false, "message": "bad payload"})
		log.Println(err)
		return
	}
	user := c.MustGet("user").(model.User)
	var group model.Group
	err = db.DB.Where("code ILIKE ?", payload.Code).Find(&group).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"status": false, "message": "invalid code"})
		log.Println(err)
		return
	}
	err = db.DB.Model(&group).Association("Users").Append(&user).Error
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"status": false, "message": "internal error"})
		log.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": true, "message": "group joined"})
}
