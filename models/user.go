package models

import (
	"gorm.io/gorm"
)//hi

type User struct {
	gorm.Model
	Name  string `json:"name"`
	Email string `json:"email" gorm:"unique"`
	Phone string `json:"phone"`
}