package repository

import (
	"bcc-project-v/src/entities"
)

func (r *Repository) SaveCart(product *entities.Cart) error {
	err := r.db.Save(&product).Error
	return err
}

func (r *Repository) GetOrCreateCart(userID uint, cart *entities.Cart) error {
	//Fungsi GORM di bawah untuk mencari cart sekaligus membuat cart apabila tidak ditemukan.
	err := r.db.Debug().FirstOrCreate(&cart, entities.Cart{UserID: userID}).Error
	return err

	// Mengambil data product di cart
	// return r.db.Preload("Products").Where("user_id = ?", userID).First(&cart).Error
}

func (r *Repository) GetProductCart(userID uint) (*entities.Cart, error) {
	cart := entities.Cart{}
	err := r.db.Preload("User").Preload("CartProducts").Where("user_id = ?", userID).First(&cart).Error
	return &cart, err
}

func (r *Repository) GetCart(userID uint) (*entities.Cart, error) {
	cart := entities.Cart{}
	err := r.db.Preload("User").Preload("Cart").Where("user_id = ?", userID).First(&cart).Error
	return &cart, err
}

func (r *Repository) DeleteCartProduct(cartID uint, productID uint) error {
	cartProduct := entities.CartProduct{}
	err := r.db.Where("cart_id = ?", cartID).Where("product_id = ?", productID).Delete(&cartProduct).Error
	return err
}

func (r *Repository) AddProductToCart(cart *entities.Cart, cartProduct *entities.CartProduct) error {
	err := r.db.Model(cart).Association("Products").Append(cartProduct)
	return err
}
