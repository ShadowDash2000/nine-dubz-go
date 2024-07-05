package usecase

import "nine-dubz/app/model"

type MovieRepository interface {
	Add(movie *model.Movie) (uint, error)
	Remove(id uint) error
	Update(movie *model.Movie) error
	Get(id uint) (*model.Movie, error)
}
