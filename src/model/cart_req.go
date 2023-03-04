package model

type NewItem struct {
	Quantity uint `json:"qty" binding:"required"`
}
