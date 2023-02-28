package entities

import "gorm.io/gorm"

type Admin struct {
	gorm.Model
	ID       uint    `json:"ID" gorm:"primaryKey" binding:"required"`
	Shop     string  `json:"shop" binding:"required"`
	Email    *string `json:"email" binding:"required"`
	Password string  `json:"password" binding:"required"`
	Product  []Product
}
