package entities

type OfflinePayment struct {
	ID          uint    `json:"ID" gorm:"primaryKey" binding:"required"`
	UserID      uint    `json:"user_id"`
	User        User    `json:"user" gorm:"foreignKey:UserID"`
	TotalPrice  float64 `json:"total_price" binding:"required"`
	PaymentCode string  `json:"payment_code" binding:"required"`
}
