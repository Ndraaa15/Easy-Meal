package entities

import (
	"time"

	"gorm.io/gorm"
)

type Cart struct {
	gorm.Model
	ID           uint          `json:"ID" gorm:"primaryKey" binding:"required"`
	UserID       uint          `json:"user_id"`
	User         User          `json:"user" gorm:"foreignKey:UserID"`
	Products     []Product     `json:"-" gorm:"many2many:cart_products"`
	CartProducts []CartProduct `json:"cart_product"`
}
type CartProduct struct {
	ID        uint `json:"ID" gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	CartID    uint    `json:"cart_id"`
	ProductID uint    `json:"product_id" binding:"required" gorm:"primaryKey"`
	Quantity  uint    `json:"qty" binding:"required"`
	Product   Product `json:"-" gorm:"foreignKey:ProductID"`
	Cart      Cart    `json:"-" gorm:"foreignKey:CartID"`
}
