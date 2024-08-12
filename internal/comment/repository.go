package comment

import (
	"gorm.io/gorm"
	"nine-dubz/internal/pagination"
)

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) Create(comment *Comment) error {
	result := r.DB.Create(comment)

	return result.Error
}

func (r *Repository) GetDistinctMultiple(where, distinct interface{}) ([]Comment, error) {
	var comments []Comment
	result := r.DB.
		Distinct(distinct).
		Where(where).
		Find(&comments)

	return comments, result.Error
}

func (r *Repository) GetMultiple(where interface{}, order string, pagination *pagination.Pagination) ([]Comment, error) {
	var comments []Comment
	result := r.DB.
		Preload("User").
		Preload("User.Picture").
		Where(where).
		Limit(pagination.Limit).
		Offset(pagination.Offset).
		Order(order).
		Find(&comments)

	return comments, result.Error
}

func (r *Repository) Count(where interface{}) (int64, error) {
	var count int64
	result := r.DB.
		Model(&Comment{}).
		Where(where).
		Count(&count)

	return count, result.Error
}

func (r *Repository) Delete(commentId, userId uint) (int64, error) {
	result := r.DB.Select("SubComments").Where("user_id = ?", userId).Delete(&Comment{ID: commentId})

	return result.RowsAffected, result.Error
}
