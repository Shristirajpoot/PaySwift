package models

import (
	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model
	UserID    uint    `json:"user_id"`
	BillID    uint    `json:"bill_id"`
	Amount    float64 `json:"amount"`
	Reference string  `json:"reference" goem:"uniqueIndex"`
	Status    string  `json:"status"`
	User      User    `json:"user" gorm:"foreignKey:UserID"`
	Bill      Bill    `json:"bill" gorm:"foreignKey:BillID"`
}