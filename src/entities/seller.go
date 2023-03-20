package entities

import "gorm.io/gorm"

type Seller struct {
	gorm.Model
	ID          uint      `json:"ID" gorm:"primaryKey"`
	Shop        string    `json:"shop" binding:"required" gorm:"unique"`
	Username    string    `json:"username" binding:"required" gorm:"unique"`
	Email       string    `json:"email" binding:"required,email" gorm:"unique"`
	Password    string    `json:"-" binding:"required"`
	Address     string    `json:"address" binding:"required"`
	City        string    `json:"city" binding:"required"`
	LinkMaps    string    `json:"link_maps" binding:"required"`
	Contact     string    `json:"contact" binding:"required,e164" gorm:"unique"`
	SellerImage string    `json:"seller_image"`
	Products    []Product `json:"-"`
}
