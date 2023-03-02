package model

type NewProduct struct {
	ImageLink   string  `json:"img_link"`
	Name        string  `json:"name" binding:"required"`
	Price       float64 `json:"price" binding:"required" `
	AdminID     uint    `json:"admin_id" gorm:"foreignKey:AdminID"`
	Description string  `json:"description" binding:"required"`
	Stock       uint    `json:"stock" binding:"required"`
}

type GetProductByID struct {
	ID uint `uri:"product_id" binding:"required"`
}

type UpdateProduct struct {
	ImageLink   string  `json:"img_link"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	Stock       uint    `json:"stock"`
}
