package role

import "nine-dubz/internal/apimethod"

type Interactor interface {
	GetApiMethodsByRolesIds(rolesIds []uint) ([]apimethod.ApiMethod, error)
}
