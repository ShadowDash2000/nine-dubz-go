package comment

import "nine-dubz/internal/pagination"

type Interactor interface {
	Create(comment *Comment) error
	GetDistinctMultiple(where, distinct interface{}) ([]Comment, error)
	GetMultiple(where interface{}, order, orderSub string, pagination, paginationSub *pagination.Pagination) ([]Comment, error)
	Count(where interface{}) (int64, error)
	Delete(commentId, userId uint) (int64, error)
}
