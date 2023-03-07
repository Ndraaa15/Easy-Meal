package model

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type SellerRegister struct {
	Shop     string `json:"shop" binding:"required" required:"true"`
	Email    string `json:"email" binding:"required" required:"true"`
	Password string `json:"password" binding:"required" required:"true"`
	Address  string `json:"address" binding:"required" required:"true"`
	Contact  string `json:"contact" binding:"required" required:"true"`
}

type SellerLogin struct {
	Email    string `json:"email" binding:"required" required:"true"`
	Password string `json:"password" binding:"required" required:"true"`
}

// type SellerUpdate struct {
// 	Shop     string
// 	Email    string
// 	Password string
// }

type GetSellerByID struct {
	ID uint `uri:"id" binding:"required"`
}

type SellerClaims struct {
	ID uint `json:"id"`
	jwt.RegisteredClaims
}

func NewAdminClaims(id uint, exp time.Duration) SellerClaims {
	return SellerClaims{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
		},
	}
}
