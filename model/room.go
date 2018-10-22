package model

import (
	"github.com/jinzhu/gorm"
)

// Room represents a room
type Room struct {
	gorm.Model
	Name    string `json:"name" gorm:"column:name"`
	GroupID int    `json:"-" gorm:"column:group_id"`
	Group   Group  `json:"-"`
}
