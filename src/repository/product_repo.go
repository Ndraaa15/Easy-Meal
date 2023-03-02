package repository

import (
	"bcc-project-v/src/entities"
)

func (r *Repository) CreateProduct(product *entities.Product) error {
	err := r.db.Create(product).Error
	return err
}

func (r *Repository) GetProductByID(idProduct uint) (*entities.Product, error) {
	product := entities.Product{}
	err := r.db.Debug().First(&product, idProduct).Error
	return &product, err
}

func (r *Repository) SaveProduct(product *entities.Product) error {
	err := r.db.Save(product).Error
	return err
}
