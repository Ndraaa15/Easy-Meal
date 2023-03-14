package entities

import "gorm.io/gorm"

type Seller struct {
	gorm.Model
	ID          uint      `json:"ID" gorm:"primaryKey" binding:"required"`
	Shop        string    `json:"shop" binding:"required" gorm:"unique"`
	Email       string    `json:"email" binding:"required,email" gorm:"unique"`
	Password    string    `json:"password" binding:"required"`
	Address     string    `json:"address" binding:"required"`
	Contact     string    `json:"contact" binding:"required,e164" gorm:"unique"`
	SellerImage string    `json:"seller_image"`
	Product     []Product `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"products"`
}
