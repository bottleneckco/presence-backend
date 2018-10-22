package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Status struct
type Status struct {
	gorm.Model
	Title     string     `json:"title" gorm:"column:title"`
	Notes     string     `json:"notes" gorm:"column:notes"`
	StartTime *time.Time `json:"start_time" gorm:"column:start_time"`
	EndTime   *time.Time `json:"end_time" gorm:"column:end_time"`
	UserID    int        `json:"-" gorm:"column:user_id"`
	User      User       `json:"-"`
}
