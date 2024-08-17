package subscription

import (
	"errors"
	"gorm.io/gorm"
	"nine-dubz/internal/pagination"
	"nine-dubz/internal/user"
)

type UseCase struct {
	SubscriptionInteractor Interactor
}

func New(db *gorm.DB) *UseCase {
	return &UseCase{
		SubscriptionInteractor: &Repository{
			DB: db,
		},
	}
}

func (uc *UseCase) Subscribe(userId, channelId uint) error {
	_, err := uc.SubscriptionInteractor.Get(userId, channelId)
	if err == nil {
		return errors.New("SUBSCRIPTION_ALREADY_EXISTS")
	}

	err = uc.SubscriptionInteractor.Create(&Subscription{
		ChannelID: channelId,
		UserID:    userId,
	})
	if err != nil {
		return errors.New("SUBSCRIPTION_FAILED_TO_SUBSCRIBE")
	}

	return nil
}

func (uc *UseCase) Unsubscribe(userId, channelId uint) error {
	rowsAffected, err := uc.SubscriptionInteractor.Delete(userId, channelId)
	if err != nil {
		return errors.New("SUBSCRIPTION_FAILED_TO_UNSUBSCRIBE")
	}
	if rowsAffected == 0 {
		return errors.New("SUBSCRIPTION_NOT_SUBSCRIBED")
	}

	return nil
}

func (uc *UseCase) Get(userId, channelId uint) (*Subscription, error) {
	return uc.SubscriptionInteractor.Get(userId, channelId)
}

func (uc *UseCase) GetAll(userId uint) ([]Subscription, error) {
	subscriptions, err := uc.SubscriptionInteractor.GetWhereMultiple(
		map[string]interface{}{"user_id": userId},
		&pagination.Pagination{
			Limit:  -1,
			Offset: -1,
		},
	)
	if err != nil {
		return nil, errors.New("SUBSCRIPTION_NO_SUBSCRIPTIONS")
	}

	return subscriptions, nil
}

func (uc *UseCase) GetMultiple(userId uint, pagination *pagination.Pagination) ([]*user.GetPublicResponse, error) {
	if pagination.Limit > 20 || pagination.Limit == -1 {
		pagination.Limit = 20
	}

	subscriptions, err := uc.SubscriptionInteractor.GetWhereMultiple(
		map[string]interface{}{"user_id": userId},
		pagination,
	)
	if err != nil {
		return nil, errors.New("SUBSCRIPTION_NO_SUBSCRIPTIONS")
	}

	var users []user.User
	for _, subscription := range subscriptions {
		users = append(users, subscription.Channel)
	}

	return user.NewGetPublicResponseMultiple(users), nil
}
