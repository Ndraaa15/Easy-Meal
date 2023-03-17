package repository

import (
	"bcc-project-v/src/entities"
	"fmt"
)

// -----------------FOR SELLER----------------------

func (r *Repository) CreateProduct(product *entities.Product) error {
	err := r.db.Debug().Create(product).Error
	return err
}

func (r *Repository) SaveProduct(product *entities.Product) error {
	err := r.db.Debug().Save(&product).Error
	return err
}

func (r *Repository) DeleteProductByID(SellerID, productID uint) error {
	product := entities.Product{}
	err := r.db.Debug().Where("seller_id = ?", SellerID).Delete(&product, productID).Error
	return err
}

func (r *Repository) GetProductByID(idProduct uint) (*entities.Product, error) {
	product := entities.Product{}
	err := r.db.Debug().Preload("Category").Preload("Seller").First(&product, idProduct).Error
	return &product, err
}

func (r *Repository) GetSellerProduct(SellerID uint) ([]entities.Product, error) {
	products := []entities.Product{}
	err := r.db.Debug().Preload("Seller").Preload("Category").Where("seller_id = ?", SellerID).Find(&products).Error
	return products, err
}

func (r *Repository) GetSellerProductByID(SellerID uint, ProductID uint) (entities.Product, error) {
	products := entities.Product{}
	err := r.db.Debug().Preload("Seller").Preload("Category").Where("seller_id = ?", SellerID).Where("ID = ?", ProductID).Find(&products).Error
	return products, err
}

// -----------------FOR BUYER----------------------

func (r *Repository) GetAllProduct() ([]entities.Product, error) {
	products := []entities.Product{}
	err := r.db.Debug().Preload("Seller").Preload("Category").Find(&products).Error
	return products, err
}

func (r *Repository) ProductsForLandingPage() ([]entities.Product, error) {
	products := []entities.Product{}
	err := r.db.Debug().Limit(3).Preload("Seller").Preload("Category").Find(&products).Error
	return products, err
}

func (r *Repository) SearchProduct(keyword string) ([]entities.Product, error) {
	fmt.Println(keyword)
	products := []entities.Product{}
	err := r.db.Debug().Preload("Seller").Preload("Category").Where("name LIKE ?", "%"+keyword+"%").Find(&products).Error
	return products, err
}

func (r *Repository) FilteredProduct(categoryID uint) ([]entities.Product, error) {
	products := []entities.Product{}
	err := r.db.Debug().Preload("Seller").Preload("Category").Where("category_id = ?", categoryID).Find(&products).Error
	return products, err
}
