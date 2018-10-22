package model

import (
	"github.com/jinzhu/gorm"
)

// User represents a User
type User struct {
	gorm.Model
	Name    string    `json:"name" gorm:"column:name"`
	Email   string    `json:"email" gorm:"column:email"`
	Picture string    `json:"picture" gorm:"column:picture"`
	Token   string    `json:"-" gorm:"token"`
	Booking []Booking `json:"-"`
	Status  []Status  `json:"-"`
	Group   []Group   `json:"-" gorm:"many2many:user_groups"`
}
