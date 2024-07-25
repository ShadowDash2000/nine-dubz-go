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

func (r *Repository) Get(commentId uint) (*Comment, error) {
	comment := &Comment{}
	result := r.DB.
		Preload("Parent").
		Preload("Parent.User").
		Preload("Parent.User.Picture").
		Preload("User").
		Preload("User.Picture").
		First(comment, commentId)

	return comment, result.Error
}

func (r *Repository) GetMultiple(where map[string]interface{}, pagination *pagination.Pagination) (*[]Comment, error) {
	var comments []Comment
	result := r.DB.
		Preload("Parent").
		Preload("Parent.User").
		Preload("Parent.User.Picture").
		Preload("User").
		Preload("User.Picture").
		Where(where).
		Limit(pagination.Limit).
		Offset(pagination.Offset).
		Find(&comments)

	return &comments, result.Error
}

func (r *Repository) Delete(commentId, userId uint) (int64, error) {
	result := r.DB.Where("user_id = ?", userId).Delete(&Comment{}, commentId)

	return result.RowsAffected, result.Error
}
