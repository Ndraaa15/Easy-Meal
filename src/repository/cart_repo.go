package repository

import (
	"bcc-project-v/src/entities"
)

func (r *Repository) AddToCart(item *entities.Cart) error {
	err := r.db.Debug().Create(item).Error
	return err
}

func (r *Repository) FindProductInCart(productID, userID uint) error {
	cart := entities.Cart{}
	err := r.db.Where("user_id = ? AND product_id = ?", userID, productID).First(&cart).Error
	return err
}

func (r *Repository) SaveCart(product *entities.Cart) error {
	err := r.db.Save(&product).Error
	return err
}

func (r *Repository) GetCartByUserID(userID uint) (*entities.Cart, error) {
	var cart entities.Cart
	if err := r.db.Preload("Products").First(&cart, "user_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &cart, nil
}

func (r *Repository) CreateCartItem(cartItem *entities.CartProduct) error {
	return r.db.Create(cartItem).Error
}

func (r *Repository) FindCartByUserID(userID uint, cart *entities.Cart) error {
	return r.db.Preload("Products.CartItem").Where("user_id = ?", userID).First(cart).Error
}

func (r *Repository) SaveToCart(cart *entities.Cart) error {
	return r.db.Omit("Products.CartItem.ID", "Products.CreatedAt", "Products.UpdatedAt", "Products.DeletedAt", "Products.Carts", "Products.CartID", "Products.ProductID").Create(cart).Error
}
