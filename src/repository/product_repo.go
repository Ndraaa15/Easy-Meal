package repository

import (
	"bcc-project-v/src/entities"
)

func (r *Repository) BindingPostAdmin(admin *entities.Admin, adminID uint) error {
	err := r.db.Preload("Posts").First(&admin, adminID).Error
	// db.Create(admin{
	// 	product: CreditCard{Number: "411111111111"}
	//   })
	// err := r.db.Create(product).Error
	// r.db.Model(product).Association("Admin").Append()
	return err
}

func (r *Repository) CreateProduct(product *entities.Product) error {
	err := r.db.Create(product).Error
	return err
}
