package controller

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"gorm.io/gorm"
	"nine-dubz/app/model"
	"time"
)

type UserRepository struct {
	DB *gorm.DB
}

func (ur *UserRepository) Add(user *model.User) error {
	role := &model.Role{}
	result := ur.DB.First(&role, "code = ?", "all")
	if result.Error != nil {
		return errors.New("user group \"all\" not found")
	}
	user.Roles = []model.Role{*role}

	if user.Password != "" {
		hash := md5.Sum([]byte(user.Password))
		user.Password = hex.EncodeToString(hash[:])
	}

	hash := md5.Sum([]byte(user.Name + user.Email + user.Password + time.Now().String()))
	user.Hash = hex.EncodeToString(hash[:])

	result = ur.DB.Create(&user)

	return result.Error
}

func (ur *UserRepository) Remove(id uint) error {
	result := ur.DB.Delete(&model.User{}, id)

	return result.Error
}

func (ur *UserRepository) Save(user *model.User) error {
	result := ur.DB.Save(&user)

	return result.Error
}

func (ur *UserRepository) Updates(user *model.User) error {
	result := ur.DB.Updates(&user)

	return result.Error
}

func (ur *UserRepository) Get(user *model.User) error {
	result := ur.DB.Preload("Picture").First(&user)

	return result.Error
}

func (ur *UserRepository) GetWhere(user *model.User, where map[string]interface{}) error {
	result := ur.DB.Preload("Picture").Where(where).First(&user)

	return result.Error
}

func (ur *UserRepository) GetById(id uint) (*model.User, error) {
	user := &model.User{}
	result := ur.DB.Preload("Picture").First(&user, id)

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
	result := ur.DB.First(&user, "active = ? AND email = ? AND password = ?", true, user.Email, user.Password)
	if result.Error != nil {
		return false
	}

	return true
}

func (ur *UserRepository) LoginWOPassword(user *model.User) bool {
	result := ur.DB.First(&user, "email = ?", user.Email)
	if result.Error != nil {
		return false
	}

	return true
}
