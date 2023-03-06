package entities

type Payment struct {
	ID         uint    `json:"ID" gorm:"primaryKey" binding:"required"`
	TotalPrice float32 `json:"total_price" binding:"required"`
	UserID     uint    `json:"user_id"`
	User       User    `json:"user" gorm:"foreignKey:UserID"`
}
