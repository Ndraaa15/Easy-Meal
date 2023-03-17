package repository

import (
	"bcc-project-v/src/entities"
)

func (r *Repository) SaveCart(cart *entities.Cart) error {
	err := r.db.Debug().Save(&cart).Error
	return err
}

func (r *Repository) GetOrCreateCart(userID uint, cart *entities.Cart) error {
	err := r.db.Debug().Preload("CartProducts").FirstOrCreate(&cart, entities.Cart{UserID: userID}).Error
	return err
}

func (r *Repository) GetProductCart(userID uint) (*entities.Cart, error) {
	cart := entities.Cart{}
	err := r.db.Debug().Preload("User").Preload("CartProducts.Product.Seller").Preload("CartProducts.Product.Category").Where("user_id = ?", userID).First(&cart).Error
	return &cart, err
}

func (r *Repository) DeleteCartProduct(cartID uint, productID uint) error {
	cartProduct := entities.CartProduct{}
	err := r.db.Debug().Where("cart_id = ?", cartID).Where("product_id = ?", productID).Delete(&cartProduct).Unscoped().Error
	return err
}

func (r *Repository) FindProduct(cartID uint, productID uint) (entities.CartProduct, error) {
	product := entities.CartProduct{}
	err := r.db.Debug().Where("cart_id = ?", cartID).Where("product_id = ?", productID).First(&product).Error
	return product, err
}

func (r *Repository) UpdateSameProduct(productCartID uint, newCartProduct *entities.CartProduct) error {
	err := r.db.Debug().Where("ID = ?", productCartID).Save(&newCartProduct).Error
	return err
}

func (r *Repository) DeleteCartProductByID(targetDeleteID uint) error {
	err := r.db.Debug().Delete(&entities.CartProduct{}, targetDeleteID).Error
	return err
}

func (r *Repository) DeleteCartProductByCartID(cartID uint) error {
	cartProduct := entities.CartProduct{}
	err := r.db.Debug().Where("cart_id = ?", cartID).Delete(&cartProduct).Unscoped().Error
	return err
}

func (r *Repository) GetProductCartForPayment(cartID uint) ([]entities.CartProduct, error) {
	cartProducts := []entities.CartProduct{}
	err := r.db.Debug().Where("cart_id = ?", cartID).Find(&cartProducts).Error
	return cartProducts, err
}

func (r *Repository) CreatePaymentProduct(productPayment *entities.PaymentProduct) error {
	err := r.db.Debug().Create(productPayment).Error
	return err
}

func (r *Repository) SavePaymentProduct(productPayment *entities.PaymentProduct) error {
	err := r.db.Debug().Save(productPayment).Error
	return err
}
