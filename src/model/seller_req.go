package model

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type SellerRegister struct {
	Shop     string `json:"shop" binding:"required"`
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Address  string `json:"address" binding:"required"`
	City     string `json:"city" binding:"required"`
	Contact  string `json:"contact" binding:"required"`
}
type SellerLogin struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password" binding:"required"`
}

type CheckOrder struct {
	PaymentCode string `json:"payment_code" binding:"required"`
}

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
