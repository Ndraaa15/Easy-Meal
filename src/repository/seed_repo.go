package repository

import (
	"bcc-project-v/src/entities"

	"gorm.io/gorm"
)

func (r *Repository) SeedCategory() error {
	var categories []entities.Category

	if err := r.db.First(&categories).Error; err != gorm.ErrRecordNotFound {
		return err
	}
	categories = []entities.Category{
		{
			Name: entities.Vegetables,
		},
		{
			Name: entities.Fruits,
		},
	}

	if err := r.db.Create(&categories).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) FindCategory(categoryProduct *entities.Category, CategoryID uint) error {
	err := r.db.Find(&categoryProduct, CategoryID).Error
	return err
}

func (r *Repository) SeedStatus() error {
	var status []entities.Status

	if err := r.db.First(&status).Error; err != gorm.ErrRecordNotFound {
		return err
	}
	status = []entities.Status{
		{
			Status: entities.Process,
		},
		{
			Status: entities.Done,
		},
		{
			Status: entities.Failed,
		},
	}

	if err := r.db.Create(&status).Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) FindStatus(statusPayment *entities.Status, StatusID uint) error {
	err := r.db.Find(&statusPayment, StatusID).Error
	return err
}
