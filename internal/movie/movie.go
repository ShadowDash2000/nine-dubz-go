package movie

import (
	"gorm.io/gorm"
	"nine-dubz/internal/file"
	"nine-dubz/model"
)

type UseCase struct {
	MovieInteractor Interactor
	FileUseCase     *file.UseCase
}

func New(db *gorm.DB, fuc *file.UseCase) *UseCase {
	return &UseCase{
		MovieInteractor: &Repository{
			DB: db,
		},
		FileUseCase: fuc,
	}
}

func (uc *UseCase) Add() (*Movie, error) {
	movie := &Movie{}
	err := uc.MovieInteractor.Add(movie)
	if err != nil {
		return nil, err
	}

	return movie, nil
}

func (uc *UseCase) Remove(movieId uint) error {
	return uc.MovieInteractor.Remove(movieId)
}

func (uc *UseCase) Updates(movie *UpdateRequest) error {
	movieRequest := NewUpdateRequest(movie)
	return uc.MovieInteractor.Updates(movieRequest)
}

func (uc *UseCase) Get(movieId uint) (*GetResponse, error) {
	movie, err := uc.MovieInteractor.Get(movieId)
	if err != nil {
		return nil, err
	}

	return NewGetResponse(movie), nil
}

func (uc *UseCase) GetAll(pagination *model.Pagination) ([]*GetResponse, error) {
	movies, err := uc.MovieInteractor.GetAll(pagination)
	if err != nil {
		return nil, err
	}

	if len(*movies) == 0 {
		return nil, err
	}

	var moviesPayload []*GetResponse
	for _, movie := range *movies {
		moviesPayload = append(moviesPayload, NewGetResponse(&movie))
	}

	return moviesPayload, nil
}
