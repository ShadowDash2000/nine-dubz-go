package comment

import "nine-dubz/internal/pagination"

type Interactor interface {
	Create(comment *Comment) error
	Get(commentId uint) (*Comment, error)
	GetMultiple(where map[string]interface{}, pagination *pagination.Pagination) (*[]Comment, error)
	Delete(commentId, userId uint) (int64, error)
}
