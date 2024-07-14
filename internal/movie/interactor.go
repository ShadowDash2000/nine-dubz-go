package movie

import (
	"nine-dubz/model"
)

type Interactor interface {
	Add(movie *Movie) error
	Remove(id uint) error
	Save(movie *Movie) error
	Updates(movie *Movie) error
	Get(id uint) (*Movie, error)
	GetAll(pagination *model.Pagination) (*[]Movie, error)
}
