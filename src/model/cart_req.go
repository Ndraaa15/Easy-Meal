package model

type ProductQuantity struct {
	Quantity uint `json:"quantity" binding:"required"`
}
