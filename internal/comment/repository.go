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

func (r *Repository) Get(where interface{}, order, orderSub string, paginationSub *pagination.Pagination) (Comment, error) {
	comment := Comment{}
	result := r.DB.
		Preload("Parent").
		Preload("User").
		Preload("User.Picture").
		Preload("SubComments", func(db *gorm.DB) *gorm.DB {
			return db.
				Order(orderSub).
				Limit(paginationSub.Limit).
				Offset(paginationSub.Offset)
		}).
		Preload("SubComments.User").
		Preload("SubComments.User.Picture").
		Where(where).
		Order(order).
		First(comment)

	return comment, result.Error
}

func (r *Repository) GetDistinctMultiple(where, distinct interface{}) ([]Comment, error) {
	var comments []Comment
	result := r.DB.
		Distinct(distinct).
		Where(where).
		Find(&comments)

	return comments, result.Error
}

func (r *Repository) GetMultiple(where interface{}, order, orderSub string, pagination, paginationSub *pagination.Pagination) ([]Comment, error) {
	var comments []Comment
	result := r.DB.
		Preload("User").
		Preload("User.Picture").
		Preload("SubComments", func(db *gorm.DB) *gorm.DB {
			return db.
				Order(orderSub).
				Limit(paginationSub.Limit).
				Offset(paginationSub.Offset)
		}).
		Preload("SubComments.User").
		Preload("SubComments.User.Picture").
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
	result := r.DB.Where("user_id = ?", userId).Delete(&Comment{}, commentId)

	return result.RowsAffected, result.Error
}
