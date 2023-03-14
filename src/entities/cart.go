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
	ID           uint `json:"ID" gorm:"primarykey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	CartID       uint    `json:"cart_id" gorm:"foreignKey:CartID"`
	ProductID    uint    `json:"product_id" binding:"required"`
	SellerID     uint    `json:"seller_id" binding:"required"`
	Quantity     uint    `json:"qty" binding:"required"`
	Seller       Seller  `json:"-" gorm:"foreignKey:SellerID"`
	Product      Product `json:"-" gorm:"foreignKey:ProductID"`
	Cart         Cart    `json:"-" gorm:"foreignKey:CartID"`
	ProductPrice float64 `json:"product_price" binding:"required"`
}

type PaymentProduct struct {
	ID           uint `json:"ID" gorm:"primarykey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	CartID       uint    `json:"cart_id" gorm:"foreignKey:CartID"`
	Cart         Cart    `json:"-" gorm:"foreignKey:CartID"`
	PaymentID    uint    `json:"payment_id" gorm:"foreignKey:PaymentID"`
	ProductID    uint    `json:"product_id" binding:"required"`
	SellerID     uint    `json:"seller_id" binding:"required"`
	Quantity     uint    `json:"qty" binding:"required"`
	Seller       Seller  `json:"-" gorm:"foreignKey:SellerID"`
	Product      Product `json:"product" gorm:"foreignKey:ProductID"`
	Payment      Payment `json:"-" gorm:"foreignKey:PaymentID"`
	ProductPrice float64 `json:"product_price" binding:"required"`
}
