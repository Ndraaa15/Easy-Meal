package model

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type RegisterUser struct {
	FName    string `json:"fname" required:"true"`
	Email    string `json:"email"  required:"true"`
	Username string `json:"username" required:"true"`
	Password string `json:"password" required:"true"`
	Address  string `json:"address" required:"true"`
	Contact  string `json:"contact" required:"true"`
}

type LoginUser struct {
	Email    string `json:"email" required:"true"`
	Username string `json:"username" required:"true"`
	Password string `json:"password" binding:"required"`
}

type GetUserByID struct {
	ID uint `uri:"id" binding:"required"`
}

type UserClaims struct {
	ID uint `json:"id"`
	jwt.RegisteredClaims
}

func NewUserClaims(id uint, exp time.Duration) UserClaims {
	return UserClaims{
		ID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
		},
	}
}
