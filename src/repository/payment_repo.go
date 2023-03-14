package repository

import "bcc-project-v/src/entities"

func (r *Repository) CreatePayment(payment *entities.Payment) error {
	err := r.db.Create(payment).Error
	return err
}

func (r *Repository) SavePayment(payment *entities.Payment) error {
	err := r.db.Save(payment).Error
	return err
}
