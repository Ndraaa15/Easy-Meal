package entities

import "gorm.io/gorm"

type Seller struct {
	gorm.Model
	ID       uint      `json:"ID" gorm:"primaryKey" binding:"required"`
	Shop     string    `json:"shop" binding:"required"`
	Email    *string   `json:"email" binding:"required"`
	Password string    `json:"password" binding:"required"`
	Product  []Product `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"products"`
}
