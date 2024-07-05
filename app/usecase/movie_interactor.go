package usecase

import "nine-dubz/app/model"

type MovieInteractor struct {
	MovieRepository MovieRepository
}

func (mi *MovieInteractor) Add(movie *model.Movie) (uint, error) {
	return mi.MovieRepository.Add(movie)
}

func (mi *MovieInteractor) Remove(id uint) error {
	return mi.MovieRepository.Remove(id)
}

func (mi *MovieInteractor) Update(movie *model.Movie) error {
	return mi.MovieRepository.Update(movie)
}

func (mi *MovieInteractor) Get(id uint) (*model.Movie, error) {
	return mi.MovieRepository.Get(id)
}
