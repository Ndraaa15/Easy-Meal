package repository

import (
	"bcc-project-v/src/entities"
)

func (r *Repository) CreateCart(cart *entities.Cart) error {
	err := r.db.Create(cart).Error
	return err
}

func (r *Repository) SaveCart(product *entities.Cart) error {
	err := r.db.Save(&product).Error
	return err
}

func (r *Repository) FindCartByUserID(userID uint, cart *entities.Cart) error {
	return r.db.Preload("Products").Where("user_id = ?", userID).First(&cart).Error
}

func (r *Repository) CreateCartProduct(product *entities.CartProduct) error {
	return r.db.Create(product).Error
}

func (r *Repository) GetCart(userID uint) (*entities.Cart, error) {
	cart := entities.Cart{}
	err := r.db.First(&cart, userID).Error
	return &cart, err
}

func (r *Repository) DeleteCartProduct(cart *entities.Cart, productID uint) error {
	err := r.db.Delete(cart, productID).Error
	return err
}
