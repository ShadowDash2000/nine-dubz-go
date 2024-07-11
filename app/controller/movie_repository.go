package controller

import (
	"gorm.io/gorm"
	"nine-dubz/app/model"
)

type MovieRepository struct {
	DB *gorm.DB
}

func NewMovieRepository(db *gorm.DB) *MovieRepository {
	return &MovieRepository{
		DB: db,
	}
}

func (mr *MovieRepository) Add(movie *model.Movie) (*model.Movie, error) {
	result := mr.DB.Create(&movie)

	return movie, result.Error
}

func (mr *MovieRepository) Remove(id uint) error {
	result := mr.DB.Delete(&model.Movie{}, id)

	return result.Error
}

func (mr *MovieRepository) Save(movie *model.Movie) error {
	result := mr.DB.Save(&movie)

	return result.Error
}

func (mr *MovieRepository) Get(id uint) (*model.Movie, error) {
	movie := &model.Movie{}
	result := mr.DB.Preload("Video").First(&movie, id)

	return movie, result.Error
}

func (mr *MovieRepository) GetAll(pagination *model.Pagination) (*[]model.Movie, error) {
	movies := &[]model.Movie{}
	result := mr.DB.
		Preload("Video").
		Limit(pagination.Limit).
		Offset(pagination.Offset).
		Where("is_published = 1").
		Find(&movies)

	return movies, result.Error
}
