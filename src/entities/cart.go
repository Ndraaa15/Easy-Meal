package entities

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	ID       uint          `json:"ID" gorm:"primaryKey" binding:"required"`
	Code     uint          `json:"code" binding:"required"`
	UserID   uint          `json:"user_id" binding:"required"`
	Products []CartProduct `gorm:"many2many:cart_products"`
}
type CartProduct struct {
	gorm.Model
	CartID    uint `json:"cart_id"`
	ProductID uint `json:"product_id"`
	Quantity  uint `json:"qty" binding:"required"`
}
