package entities

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	ID           uint          `json:"ID" gorm:"primaryKey" binding:"required"`
	UserID       uint          `json:"user_id"`
	User         User          `json:"user" gorm:"foreignKey:UserID"`
	TotalPrice   float64       `json:"total_price" binding:"required"`
	Products     []Product     `json:"-" gorm:"many2many:cart_products"`
	CartProducts []CartProduct `json:"cart_product"`
}
type CartProduct struct {
	gorm.Model
	CartID       uint    `json:"cart_id"`
	ProductID    uint    `json:"product_id" binding:"required"`
	Quantity     uint    `json:"qty" binding:"required"`
	Product      Product `json:"-" gorm:"foreignKey:ProductID"`
	Cart         Cart    `json:"-" gorm:"foreignKey:CartID"`
	ProductPrice float64 `json:"product_price" binding:"required"`
}
