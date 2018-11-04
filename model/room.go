package model

import (
	"github.com/jinzhu/gorm"
)

// Room represents a room
type Room struct {
	gorm.Model
	Name    string `json:"name" gorm:"type:text CHECK(name <> '');column:name;not null"`
	GroupID int    `json:"-" gorm:"column:group_id"`
	Group   Group  `json:"-"`
}
