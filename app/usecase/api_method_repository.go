package usecase

import "nine-dubz/app/model"

type ApiMethodRepository interface {
	Add(apiMethod *model.ApiMethod) (uint, error)
	Remove(id uint) error
	Update(apiMethod *model.ApiMethod) error
	Get(id uint) (*model.ApiMethod, error)
}
