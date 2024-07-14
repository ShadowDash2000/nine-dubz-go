package role

import (
	"gorm.io/gorm"
	"nine-dubz/internal/apimethod"
)

type UseCase struct {
	RoleInteractor Interactor
}

func New(db *gorm.DB) *UseCase {
	return &UseCase{
		RoleInteractor: &Repository{
			DB: db,
		},
	}
}

func (uc *UseCase) GetApiMethodsByRolesIds(rolesIds []uint) ([]apimethod.ApiMethod, error) {
	return uc.RoleInteractor.GetApiMethodsByRolesIds(rolesIds)
}
