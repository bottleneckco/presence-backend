package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Status struct
type Status struct {
	gorm.Model
	Title     string     `json:"title" gorm:"type:text CHECK(title <> '');column:title;not null"`
	Category  string     `json:"category" gorm:"type:text CHECK(category <> '');column:category;not null"`
	Notes     string     `json:"notes" gorm:"column:notes"`
	StartTime *time.Time `json:"start_time" gorm:"column:start_time;not null"`
	EndTime   *time.Time `json:"end_time" gorm:"column:end_time"`
	UserID    int        `json:"-" gorm:"column:user_id"`
	User      User       `json:"user"`
}
