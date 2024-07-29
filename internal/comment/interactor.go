package comment

import "nine-dubz/internal/pagination"

type Interactor interface {
	Create(comment *Comment) error
	Get(where map[string]interface{}, order, orderSub string, paginationSub *pagination.Pagination) (*Comment, error)
	GetMultiple(where map[string]interface{}, order, orderSub string, pagination, paginationSub *pagination.Pagination) (*[]Comment, error)
	Delete(commentId, userId uint) (int64, error)
}
