package entities

type Payment struct {
	ID     uint `json:"ID" gorm:"primaryKey" binding:"required"`
	UserID uint `json:"user_id"`
	User   User `json:"user" gorm:"foreignKey:UserID"`
}
