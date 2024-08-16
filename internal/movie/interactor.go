package movie

import (
	"nine-dubz/internal/pagination"
)

type Interactor interface {
	Add(movie *Movie) error
	Delete(id uint) error
	Save(movie *Movie) error
	Updates(movie *Movie) error
	UpdatesWhere(movie *Movie, where map[string]interface{}) (int64, error)
	UpdatesSelectWhere(movie *Movie, selectQuery, whereQuery interface{}) (int64, error)
	AppendAssociation(movie *Movie, association string, append interface{}) error
	Get(code string) (*Movie, error)
	GetUnscoped(code string) (*Movie, error)
	GetWhere(where interface{}) (*Movie, error)
	GetSelectWhere(selectQuery, where interface{}) (*Movie, error)
	GetWhereCount(where interface{}) (int64, error)
	GetMultipleByUserId(userId uint, pagination *pagination.Pagination, order string) (*[]Movie, error)
	GetMultiple(pagination *pagination.Pagination, order string) (*[]Movie, error)
	GetWhereMultiple(pagination *pagination.Pagination, where map[string]interface{}) (*[]Movie, error)
	GetPreloadWhere(preloads []string, whereQuery interface{}) (*Movie, error)
}
