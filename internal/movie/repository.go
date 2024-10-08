package movie

import (
	"gorm.io/gorm"
	"nine-dubz/internal/pagination"
)

type Repository struct {
	DB *gorm.DB
}

func (mr *Repository) Add(movie *Movie) error {
	result := mr.DB.Create(&movie)

	return result.Error
}

func (mr *Repository) Delete(id uint) error {
	return mr.DB.Select("Videos").Delete(&Movie{ID: id}).Error
}

func (mr *Repository) Save(movie *Movie) error {
	result := mr.DB.Save(&movie)

	return result.Error
}

func (mr *Repository) Updates(movie *Movie) error {
	result := mr.DB.Updates(&movie)

	return result.Error
}

func (mr *Repository) UpdatesWhere(movie *Movie, where map[string]interface{}) (int64, error) {
	result := mr.DB.Where(where).Updates(&movie)

	return result.RowsAffected, result.Error
}

func (mr *Repository) UpdatesSelectWhere(movie *Movie, selectQuery, whereQuery interface{}) (int64, error) {
	result := mr.DB.Select(selectQuery).Where(whereQuery).Updates(&movie)

	return result.RowsAffected, result.Error
}

func (mr *Repository) AppendAssociation(movie *Movie, association string, append interface{}) error {
	return mr.DB.Model(movie).Association(association).Append(append)
}

func (mr *Repository) Get(code string) (*Movie, error) {
	movie := &Movie{}
	result := mr.DB.
		Preload("Videos").
		Preload("Videos").
		Preload("Videos.File").
		Preload("Preview").
		Preload("PreviewWebp").
		Preload("DefaultPreview").
		Preload("DefaultPreviewWebp").
		Preload("WebVtt").
		Preload("User").
		Preload("User.Picture").
		First(&movie, "code = ?", code)

	return movie, result.Error
}

func (mr *Repository) GetWhere(where interface{}) (*Movie, error) {
	movie := &Movie{}
	result := mr.DB.
		Preload("Videos").
		Preload("Videos.File").
		Preload("Preview").
		Preload("PreviewWebp").
		Preload("DefaultPreview").
		Preload("DefaultPreviewWebp").
		Preload("WebVtt").
		Where(where).
		First(&movie)

	return movie, result.Error
}

func (mr *Repository) GetSelectWhere(selectQuery, where interface{}) (*Movie, error) {
	movie := &Movie{}
	result := mr.DB.Select(selectQuery).Where(where).First(&movie)

	return movie, result.Error
}

func (mr *Repository) GetWhereCount(where interface{}) (int64, error) {
	var count int64
	result := mr.DB.Model(&Movie{}).Where(where).Count(&count)

	return count, result.Error
}

func (mr *Repository) GetMultipleByUserId(userId uint, pagination *pagination.Pagination, order string) (*[]Movie, error) {
	movies := &[]Movie{}
	result := mr.DB.
		Preload("Videos").
		Preload("Videos.File").
		Preload("Preview").
		Preload("PreviewWebp").
		Preload("DefaultPreview").
		Preload("DefaultPreviewWebp").
		Preload("WebVtt").
		Where("user_id = ?", userId).
		Limit(pagination.Limit).
		Offset(pagination.Offset).
		Order(order).
		Find(&movies)

	return movies, result.Error
}

func (mr *Repository) GetMultiple(pagination *pagination.Pagination, order string) (*[]Movie, error) {
	movies := &[]Movie{}
	result := mr.DB.
		Preload("Videos").
		Preload("Videos.File").
		Preload("Preview").
		Preload("PreviewWebp").
		Preload("DefaultPreview").
		Preload("DefaultPreviewWebp").
		Preload("WebVtt").
		Preload("User").
		Preload("User.Picture").
		Limit(pagination.Limit).
		Offset(pagination.Offset).
		Where("is_published = 1").
		Order(order).
		Find(&movies)

	return movies, result.Error
}

func (mr *Repository) GetWhereMultiple(where map[string]interface{}, pagination *pagination.Pagination, order string) (*[]Movie, error) {
	movies := &[]Movie{}
	result := mr.DB.
		Preload("Videos").
		Preload("Videos.File").
		Preload("Preview").
		Preload("PreviewWebp").
		Preload("DefaultPreview").
		Preload("DefaultPreviewWebp").
		Preload("WebVtt").
		Preload("User").
		Preload("User.Picture").
		Limit(pagination.Limit).
		Offset(pagination.Offset).
		Where(where).
		Order(order).
		Find(&movies)

	return movies, result.Error
}

func (mr *Repository) GetPreloadWhere(preloads []string, whereQuery interface{}) (*Movie, error) {
	var movie *Movie
	result := mr.DB
	for _, preload := range preloads {
		result = result.Preload(preload)
	}

	result = result.Where(whereQuery).First(&movie)

	return movie, result.Error
}

func (mr *Repository) GetPreloadWhereMultiple(preloads []string, whereQuery interface{}, pagination *pagination.Pagination, order string) (*[]Movie, error) {
	movies := &[]Movie{}
	result := mr.DB
	for _, preload := range preloads {
		result = result.Preload(preload)
	}

	result = result.
		Limit(pagination.Limit).
		Offset(pagination.Offset).
		Where(whereQuery).
		Order(order).
		Find(&movies)

	return movies, result.Error
}
