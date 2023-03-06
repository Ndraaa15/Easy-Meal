package repository

import (
	"bcc-project-v/src/entities"
)

func (r *Repository) SaveCart(product *entities.Cart) error {
	err := r.db.Save(&product).Error
	return err
}

func (r *Repository) FindCartByUserID(userID uint, cart *entities.Cart) error {
	//Fungsi GORM di bawah untuk mencari cart sekaligus membuat cart apabila tidak ditemukan.
	err := r.db.Debug().FirstOrCreate(&cart, entities.Cart{UserID: userID}).Error
	return err

	// Mengambil data product di cart
	// return r.db.Preload("Products").Where("user_id = ?", userID).First(&cart).Error
}

func (r *Repository) GetCart(userID uint) (*entities.Cart, error) {
	cart := entities.Cart{}
	err := r.db.Where("user_id = ?", userID).First(&cart, userID).Error
	return &cart, err
}

func (r *Repository) DeleteCartProduct(cart *entities.Cart, productID uint) error {
	err := r.db.Delete(cart, productID).Error
	return err
}

func (r *Repository) AddProductToCart(cart *entities.Cart, cartProduct *entities.CartProduct) error {
	err := r.db.Model(cart).Association("CartProducts").Append(cartProduct)
	return err
}
