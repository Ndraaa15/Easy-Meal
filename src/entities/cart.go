package entities

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	ID       uint          `json:"ID" gorm:"primaryKey" binding:"required"`
	UserID   uint          `json:"user_id"`
	Products []CartProduct `gorm:"many2many:cart_products"`
}
type CartProduct struct {
	gorm.Model
	CartID    uint    `json:"cart_id"`
	ProductID uint    `json:"product_id" binding:"required"`
	Quantity  uint    `json:"qty" binding:"required"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID"`
	Cart      Cart    `json:"cart" gorm:"foreignKey:CartID"`
}
