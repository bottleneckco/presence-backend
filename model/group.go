package model

import (
	"github.com/jinzhu/gorm"
)

// Group represents a group
type Group struct {
	gorm.Model
	Name string `json:"name" gorm:"column:name"`
	Code string `json:"code" gorm:"column:code"`
	User []User `json:"-" gorm:"many2many:user_groups"`
	Room []Room `json:"-"`
}
