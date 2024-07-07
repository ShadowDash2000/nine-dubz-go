package controller

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
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
	role := &model.Role{}
	result := ur.DB.First(&role, "code = ?", "all")
	if result.Error != nil {
		return 0, errors.New("user group \"all\" not found")
	}
	user.Roles = []model.Role{*role}

	hash := md5.Sum([]byte(user.Password))
	user.Password = hex.EncodeToString(hash[:])
	result = ur.DB.Create(&user)

	return user.ID, result.Error
}

func (ur *UserRepository) Remove(id uint) error {
	result := ur.DB.Delete(&model.User{}, id)

	return result.Error
}

func (ur *UserRepository) Update(user *model.User) error {
	result := ur.DB.Save(&user)

	return result.Error
}

func (ur *UserRepository) Get(id uint) (*model.User, error) {
	user := &model.User{}
	result := ur.DB.First(&user, id)

	return user, result.Error
}

func (ur *UserRepository) GetByName(name string) (*model.User, error) {
	user := &model.User{}
	result := ur.DB.First(&user, "name = ?", name)

	return user, result.Error
}

func (ur *UserRepository) Login(user *model.User) bool {
	hash := md5.Sum([]byte(user.Password))
	user.Password = hex.EncodeToString(hash[:])
	result := ur.DB.First(&user, "email = ? AND password = ?", user.Email, user.Password)
	if result.Error != nil {
		return false
	}

	return true
}
