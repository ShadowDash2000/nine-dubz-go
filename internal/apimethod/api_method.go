package apimethod

import "gorm.io/gorm"

type UseCase struct {
	ApiMethodInteractor Interactor
}

func New(db *gorm.DB) *UseCase {
	return &UseCase{
		ApiMethodInteractor: &Repository{
			DB: db,
		},
	}
}
func (uc *UseCase) Get(path, method string) (*ApiMethod, error) {
	apiMethod := &ApiMethod{}
	err := uc.ApiMethodInteractor.GetWhere(apiMethod, map[string]interface{}{
		"path":   path,
		"method": method,
	})

	return apiMethod, err
}
