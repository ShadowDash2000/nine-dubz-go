package user

import (
	"gorm.io/gorm"
	"nine-dubz/internal/role"
)

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) Add(user *User) uint {
	roleStruct := &role.Role{}
	result := r.DB.First(&roleStruct, "code = ?", "all")
	if result.Error != nil {
		return 0
	}
	user.Roles = []role.Role{*roleStruct}

	result = r.DB.Create(&user)
	if result.Error != nil {
		return 0
	}

	return user.ID
}

func (r *Repository) Remove(id uint) error {
	result := r.DB.Delete(&User{}, id)

	return result.Error
}

func (r *Repository) Save(user *User) error {
	result := r.DB.Save(&user)

	return result.Error
}

func (r *Repository) Updates(user *User) error {
	result := r.DB.Updates(&user)

	return result.Error
}

func (r *Repository) Get(user *User) error {
	result := r.DB.Preload("Picture").First(&user)

	return result.Error
}

func (r *Repository) GetWhere(user *User, where map[string]interface{}) error {
	result := r.DB.Preload("Picture").Where(where).First(&user)

	return result.Error
}

func (r *Repository) GetById(id uint) (*User, error) {
	user := &User{}
	result := r.DB.Preload("Picture").First(&user, id)

	return user, result.Error
}

func (r *Repository) GetByName(name string) (*User, error) {
	user := &User{}
	result := r.DB.First(&user, "name = ?", name)

	return user, result.Error
}

func (r *Repository) GetRolesByUserId(userId uint) ([]role.Role, error) {
	user := &User{}
	result := r.DB.Preload("Roles").Find(&user, userId)
	if result.Error != nil {
		return nil, result.Error
	}

	return user.Roles, nil
}

func (r *Repository) Login(user *User) uint {
	result := r.DB.First(&user, "active = ? AND email = ? AND password = ?", true, user.Email, user.Password)
	if result.Error != nil {
		return 0
	}

	return user.ID
}

func (r *Repository) LoginWOPassword(user *User) uint {
	result := r.DB.First(&user, "email = ?", user.Email)
	if result.Error != nil {
		return 0
	}

	return user.ID
}
