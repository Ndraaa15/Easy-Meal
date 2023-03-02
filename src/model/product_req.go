package model

type NewProduct struct {
	ImageLink   string  `json:"img_link" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Price       float64 `json:"price" binding:"required" `
	AdminID     uint    `json:"admin_id" gorm:"foreignKey:AdminID"`
	Description string  `json:"description" binding:"required"`
	Stock       uint    `json:"stock" binding:"required"`
}
