package subscription

import (
	"gorm.io/gorm"
	"nine-dubz/internal/pagination"
)

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) Create(sub *Subscription) error {
	return r.DB.Create(sub).Error
}

func (r *Repository) Delete(userId, channelId uint) (int64, error) {
	result := r.DB.Delete(&Subscription{}, "user_id = ? AND channel_id = ?", userId, channelId)

	return result.RowsAffected, result.Error
}

func (r *Repository) Get(userId, channelId uint) (*Subscription, error) {
	var subscription Subscription
	result := r.DB.First(&subscription, "user_id = ? AND channel_id = ?", userId, channelId)

	return &subscription, result.Error
}

func (r *Repository) GetWhereMultiple(where interface{}, pagination *pagination.Pagination) ([]Subscription, error) {
	var subscriptions []Subscription
	result := r.DB.
		Preload("Channel").
		Preload("Channel.Picture").
		Where(where).
		Limit(pagination.Limit).
		Offset(pagination.Offset).
		Find(&subscriptions)

	return subscriptions, result.Error
}
