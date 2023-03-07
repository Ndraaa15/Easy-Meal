package model

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
