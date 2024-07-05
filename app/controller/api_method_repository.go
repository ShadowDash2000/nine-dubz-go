package controller

import (
	"gorm.io/gorm"
	"nine-dubz/app/model"
)

type ApiMethodRepository struct {
	DB *gorm.DB
}

func NewApiMethodRepository(db *gorm.DB) *ApiMethodRepository {
	return &ApiMethodRepository{
		DB: db,
	}
}

func (amr *ApiMethodRepository) Add(apiMethod *model.ApiMethod) (uint, error) {
	result := amr.DB.Create(apiMethod)

	return apiMethod.ID, result.Error
}

func (amr *ApiMethodRepository) Remove(id uint) error {
	result := amr.DB.Delete(&model.ApiMethod{}, id)

	return result.Error
}

func (amr *ApiMethodRepository) Update(apiMethod *model.ApiMethod) error {
	result := amr.DB.Save(apiMethod)

	return result.Error
}

func (amr *ApiMethodRepository) Get(id uint) (*model.ApiMethod, error) {
	apiMethod := &model.ApiMethod{}
	result := amr.DB.First(apiMethod, id)

	return apiMethod, result.Error
}
