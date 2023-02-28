package repository

import (
	"bcc-project-v/src/entities"
	"bcc-project-v/src/model"
)

func (r *Repository) CreateAdmin(seller *entities.Admin) error {
	if err := r.db.Debug().Create(seller).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) FindAdminByEmail(model *model.AdminLogin) (*entities.Admin, error) {
	adminFound := entities.Admin{}
	err := r.db.Debug().Where("email = ?", model.Email).First(&adminFound).Error
	return &adminFound, err
}

func (r *Repository) FindAdminByID(ID uint) (*entities.Admin, error) {
	admin := entities.Admin{}
	err := r.db.Debug().First(&admin, ID).Error
	return &admin, err
}

func (r *Repository) UpdateAdmin(admin *entities.Admin) error {
	err := r.db.Debug().Save(&admin).Error
	return err
}
