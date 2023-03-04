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
