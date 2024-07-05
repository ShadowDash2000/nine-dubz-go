package controller

import (
	"gorm.io/gorm"
	"nine-dubz/app/model"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		DB: db,
	}
}

func (ur *UserRepository) Add(user *model.User) (uint, error) {
	result := ur.DB.Create(user)

	return user.ID, result.Error
}

func (ur *UserRepository) Remove(id uint) error {
	result := ur.DB.Delete(&model.User{}, id)

	return result.Error
}

func (ur *UserRepository) Update(user *model.User) error {
	result := ur.DB.Save(user)

	return result.Error
}

func (ur *UserRepository) Get(id uint) (*model.User, error) {
	user := &model.User{}
	result := ur.DB.First(user, id)

	return user, result.Error
}
