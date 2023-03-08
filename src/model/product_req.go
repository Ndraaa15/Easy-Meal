package model

type GetProductByID struct {
	ID uint `uri:"product_id" binding:"required"`
}
