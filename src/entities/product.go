package entities

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	ID           uint          `json:"ID" gorm:"primaryKey" binding:"required"`
	ProductImage string        `json:"product_image" binding:"required"`
	Name         string        `json:"name" binding:"required"`
	Price        float64       `json:"price" binding:"required"`
	SellerID     uint          `json:"seller_id" binding:"required"`
	Seller       Seller        `json:"Seller" gorm:"foreignKey:SellerID"`
	Description  string        `json:"description" binding:"required"`
	Stock        uint          `json:"stock" binding:"required"`
	Mass         uint          `json:"mass" binding:"required"`
	CategoryID   uint          `json:"category_id" binding:"required"`
	Category     Category      `gorm:"foreignKey:CategoryID"`
	Cart         []Cart        `json:"-" gorm:"many2many:cart_products"`
	CartProducts []CartProduct `json:"-"`
}
