package entities

import "gorm.io/gorm"

type Payment struct {
	gorm.Model
	ID          uint    `json:"ID" gorm:"primaryKey" binding:"required"`
	UserID      uint    `json:"user_id"`
	CartID      uint    `json:"cart_id"`
	User        User    `json:"-" gorm:"foreignKey:UserID"`
	Cart        Cart    `json:"-" gorm:"foreignKey:CartID"`
	Type        string  `json:"type" binding:"required"`
	TotalPrice  float64 `json:"total_price" binding:"required"`
	PaymentCode string  `json:"payment_code" binding:"required"`
	StatusID    uint    `json:"status_id" binding:"required"`
	Status      Status  `gorm:"foreignKey:StatusID"`
	FName       string  `json:"fname" binding:"required"`
	Contact     string  `json:"contact" binding:"required"`
	Address     string  `json:"address" binding:"required"`
	Email       string  `json:"email" binding:"required"`
	City        string  `json:"city" binding:"required"`
}
