package repository

import "bcc-project-v/src/entities"

func (r *Repository) GetHistory(userID uint) ([]entities.Payment, error) {
	history := []entities.Payment{}
	err := r.db.Debug().Preload("Status").Where("user_id = ?", userID).Find(&history).Error
	return history, err
}
