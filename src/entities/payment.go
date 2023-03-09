package entities

type OfflinePayment struct {
	ID          uint    `json:"ID" gorm:"primaryKey" binding:"required"`
	UserID      uint    `json:"user_id"`
	CartID      uint    `json:"cart_id"`
	User        User    `json:"-" gorm:"foreignKey:UserID"`
	Cart        Cart    `json:"-" gorm:"foreignKey:CartID"`
	TotalPrice  float64 `json:"total_price" binding:"required"`
	PaymentCode string  `json:"payment_code" binding:"required"`
}
