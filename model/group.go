package model

import (
	"github.com/jinzhu/gorm"
)

// Group represents a group
type Group struct {
	gorm.Model
	Name     string `json:"name" gorm:"column:name;not null"`
	Code     string `json:"code" gorm:"column:code;not null"`
	Users    []User `json:"-" gorm:"many2many:user_groups"`
	Rooms    []Room `json:"-"`
	Author   User   `json:"author"`
	AuthorID int    `json:"-" gorm:"column:author_id;not null"`
}
