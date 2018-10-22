package model

import (
	"github.com/jinzhu/gorm"
)

// Booking represents a booking
type Booking struct {
	gorm.Model
	Title  string `json:"title" gorm:"column:title"`
	UserID int    `json:"-" gorm:"column:user_id"`
	User   User   `json:"-"`
	Room   Room   `json:"-"`
	RoomID int    `json:"-" gorm:"column:room_id"`
}
