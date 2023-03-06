package entities

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	ID          uint    `json:"ID" gorm:"primaryKey" binding:"required"`
	ImageLink   string  `json:"img_link" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Price       float64 `json:"price" binding:"required"`
	SellerID    uint    `json:"seller_id" binding:"required" gorm:"foreignKey:SellerID"`
	Description string  `json:"description" binding:"required"`
	Stock       uint    `json:"stock" binding:"required"`
	Carts       []Cart  `json:"carts" gorm:"many2many:product_cart;"`
}
