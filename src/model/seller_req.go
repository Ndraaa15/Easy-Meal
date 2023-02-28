package model

type RegisterSeller struct {
	Shop     string `json:"shop" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Address  string `json:"address" binding:"required"`
	Contact  string `json:"contact"  binding:"required"`
}

type LoginSeller struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateSeller struct {
	Shop     string
	Email    string
	Password string
	Address  string
	Contact  string
}

type GetSellerByID struct {
	ID uint `uri:"id" binding:"required"`
}
