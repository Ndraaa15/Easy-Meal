package repository

import (
	"bcc-project-v/src/entities"
)

func (r *Repository) SaveCart(cart *entities.Cart) error {
	err := r.db.Save(&cart).Error
	return err
}

func (r *Repository) GetOrCreateCart(userID uint, cart *entities.Cart) error {
	//Fungsi GORM di bawah untuk mencari cart sekaligus membuat cart apabila tidak ditemukan.
	err := r.db.Debug().Preload("CartProducts").FirstOrCreate(&cart, entities.Cart{UserID: userID}).Error
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
func (r *Repository) GetCartForPayment(userID uint) (*entities.Cart, error) {
	cart := entities.Cart{}
	err := r.db.Where("user_id = ?", userID).First(&cart).Error
	return &cart, err
}

func (r *Repository) DeleteCartProduct(cartID uint, productID uint) error {
	cartProduct := entities.CartProduct{}
	err := r.db.Where("cart_id = ?", cartID).Where("product_id = ?", productID).Delete(&cartProduct).Error
	return err
}

func (r *Repository) FindProduct(cartID uint, productID uint) (entities.CartProduct, error) {
	product := entities.CartProduct{}
	err := r.db.Where("cart_id = ?", cartID).Where("product_id = ?", productID).First(&product).Error
	return product, err
}

// func (r *Repository) ReplaceCartProduct(cart *entities.Cart, newCartProduct *entities.CartProduct) error {
// 	err := r.db.Model(cart).Association("CartProducts").Replace(entities.Cart{CartProducts: newCartProduct})
// 	return err
// }

func (r *Repository) UpdateSameProduct(productCartID uint, newCartProduct *entities.CartProduct) error {
	err := r.db.Where("ID = ?", productCartID).Save(&newCartProduct).Error
	return err
}
