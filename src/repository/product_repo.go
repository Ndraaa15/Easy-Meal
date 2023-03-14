package repository

import (
	"bcc-project-v/src/entities"
)

// -----------------FOR SELLER----------------------

func (r *Repository) CreateProduct(product *entities.Product) error {
	err := r.db.Create(product).Error
	return err
}

func (r *Repository) SaveProduct(product *entities.Product) error {
	err := r.db.Save(&product).Error
	return err
}

func (r *Repository) DeleteProductByID(SellerID, productID uint) error {
	product := entities.Product{}
	err := r.db.Where("seller_id = ?", SellerID).Delete(&product, productID).Error
	return err
}

func (r *Repository) GetProductByID(idProduct uint) (*entities.Product, error) {
	product := entities.Product{}
	err := r.db.Debug().First(&product, idProduct).Error
	return &product, err
}

func (r *Repository) GetSellerProduct(SellerID uint) ([]entities.Product, error) {
	products := []entities.Product{}
	err := r.db.Preload("Category").Where("seller_id = ?", SellerID).Find(&products).Error
	return products, err
}

func (r *Repository) GetSellerProductByID(SellerID uint, ProductID uint) (entities.Product, error) {
	products := entities.Product{}
	err := r.db.Preload("Category").Where("seller_id = ?", SellerID).Where("ID = ?", ProductID).Find(&products).Error
	return products, err
}

// -----------------FOR BUYER----------------------

func (r *Repository) GetAllProduct() ([]entities.Product, error) {
	products := []entities.Product{}
	err := r.db.Preload("Seller").Preload("Category").Find(&products).Error
	return products, err
}

func (r *Repository) SearchProduct(keyword string) ([]entities.Product, error) {
	products := []entities.Product{}
	err := r.db.Preload("Category").Where("name like ?", "%"+keyword+"%").Find(&products).Error
	return products, err
}

func (r *Repository) FilteredProduct(categoryID uint) ([]entities.Product, error) {
	products := []entities.Product{}
	err := r.db.Preload("Category").Where("category_id = ?", categoryID).Find(&products).Error
	return products, err
}
