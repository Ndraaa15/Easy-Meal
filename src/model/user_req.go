package model

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type RegisterUser struct {
	FName    string `json:"fname" binding:"required" required:"true"`
	Email    string `json:"email" binding:"required" required:"true"`
	Username string `json:"username" binding:"required" required:"true"`
	Gender   string `json:"gender" binding:"required" required:"true"`
	Password string `json:"password" binding:"required" required:"true"`
	Address  string `json:"address" binding:"required" required:"true"`
	Contact  string `json:"contact" binding:"required" required:"true"`
}

type LoginUser struct {
	Email    string `json:"email" required:"true"`
	Username string `json:"username" required:"true"`
	Password string `json:"password" binding:"required" required:"true"`
}

type UpdateUser struct {
	FName    string `json:"fname"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Gender   string `json:"gender"`
	Password string `json:"password"`
	Address  string `json:"address"`
	Contact  string `json:"contact"`
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
