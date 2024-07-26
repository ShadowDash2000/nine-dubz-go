package movie

import (
	"nine-dubz/internal/pagination"
)

type Interactor interface {
	Add(movie *Movie) error
	Delete(userId uint, code string) error
	Save(movie *Movie) error
	Updates(movie *Movie) error
	UpdatesWhere(movie *Movie, where map[string]interface{}) (int64, error)
	Get(code string) (*Movie, error)
	GetUnscoped(code string) (*Movie, error)
	GetWhere(code string, where map[string]interface{}) (*Movie, error)
	GetMultipleByUserId(userId uint, pagination *pagination.Pagination) (*[]Movie, error)
	GetMultiple(pagination *pagination.Pagination) (*[]Movie, error)
	GetWhereMultiple(pagination *pagination.Pagination, where map[string]interface{}) (*[]Movie, error)
}
