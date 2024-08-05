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

func (mr *Repository) Delete(userId uint, code string) error {
	result := mr.DB.Where("user_id = ? AND code = ?", userId, code).Delete(&Movie{})

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

func (mr *Repository) UpdatesWhere(movie *Movie, where map[string]interface{}) (int64, error) {
	result := mr.DB.Where(where).Updates(&movie)

	return result.RowsAffected, result.Error
}

func (mr *Repository) UpdatesSelectWhere(movie *Movie, selectQuery, whereQuery interface{}) (int64, error) {
	result := mr.DB.Select(selectQuery).Where(whereQuery).Updates(&movie)

	return result.RowsAffected, result.Error
}

func (mr *Repository) Get(code string) (*Movie, error) {
	movie := &Movie{}
	result := mr.DB.
		Preload("VideoTmp").
		Preload("VideoTmp.File").
		Preload("Video").
		Preload("Video.File").
		Preload("VideoShakal").
		Preload("VideoShakal.File").
		Preload("Video360").
		Preload("Video360.File").
		Preload("Video480").
		Preload("Video480.File").
		Preload("Video720").
		Preload("Video720.File").
		Preload("Preview").
		Preload("DefaultPreview").
		Preload("WebVtt").
		Preload("User").
		Preload("User.Picture").
		First(&movie, "code = ?", code)

	return movie, result.Error
}

func (mr *Repository) GetUnscoped(code string) (*Movie, error) {
	movie := &Movie{}
	result := mr.DB.
		Unscoped().
		Preload("Video").
		Preload("Video.File").
		Preload("VideoShakal").
		Preload("VideoShakal.File").
		Preload("Video360").
		Preload("Video360.File").
		Preload("Video480").
		Preload("Video480.File").
		Preload("Video720").
		Preload("Video720.File").
		Preload("Preview").
		Preload("DefaultPreview").
		Preload("WebVtt").
		First(&movie, "code = ?", code)

	return movie, result.Error
}

func (mr *Repository) GetWhere(code string, where map[string]interface{}) (*Movie, error) {
	movie := &Movie{}
	result := mr.DB.
		Preload("Video").
		Preload("Video.File").
		Preload("VideoShakal").
		Preload("VideoShakal.File").
		Preload("Video360").
		Preload("Video360.File").
		Preload("Video480").
		Preload("Video480.File").
		Preload("Video720").
		Preload("Video720.File").
		Preload("Preview").
		Preload("DefaultPreview").
		Preload("WebVtt").
		Where(where).
		First(&movie, "code = ?", code)

	return movie, result.Error
}

func (mr *Repository) GetMultipleByUserId(userId uint, pagination *pagination.Pagination) (*[]Movie, error) {
	movies := &[]Movie{}
	result := mr.DB.
		Preload("Video").
		Preload("Video.File").
		Preload("VideoShakal").
		Preload("VideoShakal.File").
		Preload("Video360").
		Preload("Video360.File").
		Preload("Video480").
		Preload("Video480.File").
		Preload("Video720").
		Preload("Video720.File").
		Preload("Preview").
		Preload("DefaultPreview").
		Preload("WebVtt").
		Where("user_id = ?", userId).
		Limit(pagination.Limit).
		Offset(pagination.Offset).
		Find(&movies)

	return movies, result.Error
}

func (mr *Repository) GetMultiple(pagination *pagination.Pagination, order string) (*[]Movie, error) {
	movies := &[]Movie{}
	result := mr.DB.
		Preload("Video").
		Preload("Video.File").
		Preload("VideoShakal").
		Preload("VideoShakal.File").
		Preload("Video360").
		Preload("Video360.File").
		Preload("Video480").
		Preload("Video480.File").
		Preload("Video720").
		Preload("Video720.File").
		Preload("Preview").
		Preload("DefaultPreview").
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

func (mr *Repository) GetWhereMultiple(pagination *pagination.Pagination, where map[string]interface{}) (*[]Movie, error) {
	movies := &[]Movie{}
	result := mr.DB.
		Preload("VideoTmp").
		Preload("VideoTmp.File").
		Preload("Video").
		Preload("Video.File").
		Preload("VideoShakal").
		Preload("VideoShakal.File").
		Preload("Video360").
		Preload("Video360.File").
		Preload("Video480").
		Preload("Video480.File").
		Preload("Video720").
		Preload("Video720.File").
		Preload("Preview").
		Preload("DefaultPreview").
		Preload("WebVtt").
		Preload("User").
		Preload("User.Picture").
		Limit(pagination.Limit).
		Offset(pagination.Offset).
		Where(where).
		Find(&movies)

	return movies, result.Error
}
