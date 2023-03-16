package repository

import "bcc-project-v/src/entities"

func (r *Repository) CreatePayment(payment *entities.Payment) error {
	err := r.db.Debug().Preload("Products").Create(payment).Error
	return err
}

func (r *Repository) SavePayment(payment *entities.Payment) error {
	err := r.db.Debug().Save(payment).Error
	return err
}

func (r *Repository) FilteredStatus(statusID uint) ([]entities.Payment, error) {
	payments := []entities.Payment{}
	err := r.db.Debug().Where("status_id = ?", statusID).Find(&payments).Error
	return payments, err
}

func (r *Repository) GetOrder(sellerID uint) ([]entities.PaymentProduct, error) {
	productsOrder := []entities.PaymentProduct{}
	err := r.db.Debug().Where("seller_id = ?", sellerID).Find(&productsOrder).Error
	return productsOrder, err
}
