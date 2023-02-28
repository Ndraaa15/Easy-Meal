package model

type AdminRegister struct {
	Shop     string  `json:"shop" binding:"required"`
	Email    *string `json:"email" binding:"required"`
	Password string  `json:"password" binding:"required"`
}

type AdminLogin struct {
	Email    *string `json:"email" binding:"required"`
	Password string  `json:"password" binding:"required"`
}

type AdminUpdate struct {
	Shop     string
	Email    *string
	Password string
}

type GetAdminByID struct {
	ID uint `uri:"id" binding:"required"`
}
