package users

import (
	"gorm.io/gorm"
)

type Repository interface {
	Save(user User) (User, error)
	FindByEmail(email string) (User, error)
	FindByUsername(username string) (User, error)
	FindByID(ID int) (User, error)
	FindByToken(token string) (User, error)
	FindByUUID(UUID string) (User, error)
	FindByIdAndUpdateToken(ID int, token string) (User, error)
	FindByIdAndUpdateUUID(ID int, generateUUID string) (User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Save(user User) (User, error) {
	err := r.db.Create(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *repository) FindByEmail(email string) (User, error) {
	var user User
	err := r.db.Where("email = ?", email).Find(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *repository) FindByUsername(username string) (User, error) {
	var user User
	err := r.db.Where("username = ?", username).Find(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *repository) FindByID(ID int) (User, error) {
	var user User
	err := r.db.Where("id = ?", ID).Find(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *repository) FindByToken(token string) (User, error) {
	var user User
	err := r.db.Where("token = ?", token).Find(&user).Error

	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *repository) FindByIdAndUpdateToken(ID int, token string) (User, error) {
	var user User
	err := r.db.Model(&user).Where("id = ?", ID).Update("token", token).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *repository) FindByIdAndUpdateUUID(ID int, generateUUID string) (User, error) {
	var user User
	err := r.db.Model(&user).Where("id = ?", ID).Update("uuid", generateUUID).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *repository) FindByUUID(UUID string) (User, error) {
	var user User
	err := r.db.Where("uuid = ?", UUID).Find(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}
