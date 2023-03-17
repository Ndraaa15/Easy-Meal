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
	err := r.db.Debug().Model(&entities.Payment{}).Preload("PaymentProducts.Product.Seller").Preload("PaymentProducts.Product.Category").Where("status_id = ?", statusID).Find(&payments).Error
	return payments, err
}

func (r *Repository) GetOrder(sellerID uint) ([]entities.PaymentProduct, error) {
	productsOrder := []entities.PaymentProduct{}
	err := r.db.Debug().Preload("Product").Where("seller_id = ?", sellerID).Find(&productsOrder).Error
	return productsOrder, err
}

func (r *Repository) GetPaymentProductByID(productID uint) (*entities.PaymentProduct, error) {
	product := entities.PaymentProduct{}
	err := r.db.Debug().First(&product, productID).Error
	return &product, err
}
