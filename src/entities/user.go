package entities

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID        uint   `json:"ID" binding:"required" gorm:"primaryKey"`
	FName     string `json:"fname" binding:"required"`
	Email     string `json:"email" binding:"required,email" gorm:"unique"`
	Username  string `json:"username" binding:"required" gorm:"unique"`
	Gender    string `json:"gender" binding:"required"`
	Password  string `json:"-" binding:"required"`
	Address   string `json:"address" binding:"required"`
	Contact   string `json:"contact" binding:"required,e164" gorm:"unique"`
	UserImage string `json:"user_image"`
}
