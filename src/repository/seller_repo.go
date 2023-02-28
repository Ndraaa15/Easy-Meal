package repository

import (
	"bcc-project-v/src/entities"
	"bcc-project-v/src/model"
)

func (r *Repository) CreateSeller(seller *entities.Seller) error {
	err := r.db.Create(&seller).Error
	if err != nil {
		return err
	}
	return err
}

func (r *Repository) FindSellerByEmail(seller model.LoginSeller) (*entities.Seller, error) {
	sellerFound := entities.Seller{}
	err := r.db.Where("username = ?", seller.Email).First(&seller).Error
	return &sellerFound, err
}

func (r *Repository) FindSellerByID(ID uint) (*entities.Seller, error) {
	seller := entities.Seller{}
	err := r.db.First(&seller, ID).Error
	return &seller, err
}

func (r *Repository) UpdateSeller(seller *entities.Seller) error {
	err := r.db.Save(&seller).Error
	return err
}
