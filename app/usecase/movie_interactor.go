package usecase

import "nine-dubz/app/model"

type MovieInteractor struct {
	MovieRepository MovieRepository
}

func (mi *MovieInteractor) Add(movie *model.Movie) error {
	return mi.MovieRepository.Add(movie)
}

func (mi *MovieInteractor) Remove(id uint) error {
	return mi.MovieRepository.Remove(id)
}

func (mi *MovieInteractor) Save(movie *model.Movie) error {
	return mi.MovieRepository.Save(movie)
}

func (mi *MovieInteractor) Updates(movie *model.Movie) error {
	return mi.MovieRepository.Updates(movie)
}

func (mi *MovieInteractor) Get(id uint) (*model.Movie, error) {
	return mi.MovieRepository.Get(id)
}

func (mi *MovieInteractor) GetAll(pagination *model.Pagination) (*[]model.Movie, error) {
	return mi.MovieRepository.GetAll(pagination)
}
