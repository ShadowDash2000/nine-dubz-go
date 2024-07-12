package usecase

import "nine-dubz/app/model"

type MovieRepository interface {
	Add(movie *model.Movie) error
	Remove(id uint) error
	Save(movie *model.Movie) error
	Updates(movie *model.Movie) error
	Get(id uint) (*model.Movie, error)
	GetAll(pagination *model.Pagination) (*[]model.Movie, error)
}
