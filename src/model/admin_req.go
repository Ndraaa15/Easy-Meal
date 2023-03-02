package model

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type AdminRegister struct {
	Shop     string  `json:"shop" binding:"required" required:"true"`
	Email    *string `json:"email" binding:"required" required:"true"`
	Password string  `json:"password" binding:"required" required:"true"`
}

type AdminLogin struct {
	Email    *string `json:"email" binding:"required" required:"true"`
	Password string  `json:"password" binding:"required" required:"true"`
}

type AdminUpdate struct {
	Shop     string
	Email    *string
	Password string
}

type GetAdminByID struct {
	ID uint `uri:"id" binding:"required"`
}

type AdminClaims struct {
	ID uint `json:"id"`
	jwt.RegisteredClaims
}

func NewAdminClaims(id uint, exp time.Duration) AdminClaims {
	return AdminClaims{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
		},
	}
}
