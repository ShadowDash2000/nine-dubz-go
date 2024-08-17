package subscription

import "nine-dubz/internal/pagination"

type Interactor interface {
	Create(sub *Subscription) error
	Delete(userId, channelId uint) (int64, error)
	Get(userId, channelId uint) (*Subscription, error)
	GetWhereMultiple(where interface{}, pagination *pagination.Pagination) ([]Subscription, error)
}
