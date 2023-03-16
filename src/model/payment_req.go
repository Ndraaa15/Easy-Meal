package model

type DataBuyer struct {
	FName   string `json:"fname" binding:"required"`
	Contact string `json:"contact" binding:"required"`
	Address string `json:"address" binding:"required"`
	Email   string `json:"email" binding:"required"`
	City    string `json:"city" binding:"required"`
}
