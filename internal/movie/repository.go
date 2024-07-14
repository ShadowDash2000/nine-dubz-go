package movie

import (
	"gorm.io/gorm"
	"nine-dubz/model"
)

type Repository struct {
	DB *gorm.DB
}

func (mr *Repository) Add(movie *Movie) error {
	result := mr.DB.Create(&movie)

	return result.Error
}

func (mr *Repository) Remove(id uint) error {
	result := mr.DB.Delete(&Movie{}, id)

	return result.Error
}

func (mr *Repository) Save(movie *Movie) error {
	result := mr.DB.Save(&movie)

	return result.Error
}

func (mr *Repository) Updates(movie *Movie) error {
	result := mr.DB.Updates(&movie)

	return result.Error
}

func (mr *Repository) Get(id uint) (*Movie, error) {
	movie := &Movie{}
	result := mr.DB.Preload("Video").First(&movie, id)

	return movie, result.Error
}

func (mr *Repository) GetAll(pagination *model.Pagination) (*[]Movie, error) {
	movies := &[]Movie{}
	result := mr.DB.
		Preload("Video").
		Limit(pagination.Limit).
		Offset(pagination.Offset).
		Where("is_published = 1").
		Find(&movies)

	return movies, result.Error
}
