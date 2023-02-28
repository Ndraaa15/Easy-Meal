package entities

import "gorm.io/gorm"

type Wishlist struct {
	gorm.Model
	ID         uint
	User_ID    uint
	Product_ID uint
	Quantity   uint
}
