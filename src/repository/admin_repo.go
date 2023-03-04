package repository

import (
	"bcc-project-v/src/entities"
	"bcc-project-v/src/model"
)

func (r *Repository) CreateAdmin(seller *entities.Seller) error {
	if err := r.db.Debug().Create(seller).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) FindAdminByEmail(model *model.SellerLogin) (*entities.Seller, error) {
	adminFound := entities.Seller{}
	err := r.db.Debug().Where("email = ?", model.Email).First(&adminFound).Error
	return &adminFound, err
}

func (r *Repository) FindAdminByID(ID uint) (*entities.Seller, error) {
	admin := entities.Seller{}
	err := r.db.Debug().First(&admin, ID).Error
	return &admin, err
}

func (r *Repository) UpdateAdmin(admin *entities.Seller) error {
	err := r.db.Debug().Save(&admin).Error
	return err
}
