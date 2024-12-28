package models

import (
	"time"

	"gorm.io/gorm"
)

type Bill struct {
	gorm.Model
	UserID      uint      `json:"user_id"`
	Amount      float64   `json:"amount"`
	DueDate     time.Time `json:"due_date"`
	Status      string    `json:"status"`
	Description string    `json:"description"`
	User        User      `json:"user" gorm:"foreignKey:UserID"`
}