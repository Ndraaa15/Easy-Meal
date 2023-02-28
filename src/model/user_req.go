package model

type RegisterUser struct {
	FName    string `json:"fname" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Username string `json:"username" binding:"required"`
	Gender   bool   `json:"gender" binding:"required"`
	Password string `json:"password" binding:"required"`
	Address  string `json:"address" binding:"required"`
	Contact  string `json:"contact" binding:"required"`
}

type LoginUser struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password" binding:"required"`
}

type UpdateUser struct {
	FName    string `json:"fname"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Gender   bool   `json:"gender"`
	Password string `json:"password"`
	Address  string `json:"address"`
	Contact  string `json:"contact"`
}

type GetUserByID struct {
	ID uint `uri:"id" binding:"required"`
}
