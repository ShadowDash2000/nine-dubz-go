package usecase

import "nine-dubz/app/model"

type ApiMethodInteractor struct {
	ApiMethodRepository ApiMethodRepository
}

func (ami *ApiMethodInteractor) Add(apiMethod *model.ApiMethod) (uint, error) {
	return ami.ApiMethodRepository.Add(apiMethod)
}

func (ami *ApiMethodInteractor) Remove(id uint) error {
	return ami.ApiMethodRepository.Remove(id)
}

func (ami *ApiMethodInteractor) Update(apiMethod *model.ApiMethod) error {
	return ami.ApiMethodRepository.Update(apiMethod)
}

func (ami *ApiMethodInteractor) Get(id uint) (*model.ApiMethod, error) {
	return ami.ApiMethodRepository.Get(id)
}
