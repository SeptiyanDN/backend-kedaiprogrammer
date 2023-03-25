package businesses

import "gorm.io/gorm"

type Repository interface {
	FindAll() ([]Business, error)
	Save(business Business) (Business, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Save(business Business) (Business, error) {
	err := r.db.Create(&business).Error
	if err != nil {
		return business, err
	}
	return business, nil
}

func (r *repository) FindAll() ([]Business, error) {
	var businesses []Business
	err := r.db.Find(&businesses).Error
	if err != nil {
		return businesses, err
	}
	return businesses, nil
}
