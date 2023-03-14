package repository

import (
	"bcc-project-v/src/entities"
	"bcc-project-v/src/model"
)

func (r *Repository) CreateSeller(seller *entities.Seller) error {
	if err := r.db.Debug().Create(seller).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) FindSellerByEmail(model *model.SellerLogin) (*entities.Seller, error) {
	sellerFound := entities.Seller{}
	err := r.db.Debug().Where("email = ?", model.Email).First(&sellerFound).Error
	return &sellerFound, err
}

func (r *Repository) FindSellerByID(ID uint) (*entities.Seller, error) {
	seller := entities.Seller{}
	err := r.db.Debug().First(&seller, ID).Error
	return &seller, err
}

func (r *Repository) UpdateSeller(seller *entities.Seller) error {
	err := r.db.Debug().Save(&seller).Error
	return err
}

func (r *Repository) CheckOrder(paymentCode string) (*entities.Payment, error) {
	payment := entities.Payment{}
	err := r.db.Where("payment_code = ?", paymentCode).Find(&payment).Error
	return &payment, err
}
