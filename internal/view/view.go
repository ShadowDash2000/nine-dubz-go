package view

import (
	"errors"
	"gorm.io/gorm"
	"net"
	"time"
)

type UseCase struct {
	ViewInteractor Interactor
}

func New(db *gorm.DB) *UseCase {
	return &UseCase{
		ViewInteractor: &Repository{
			DB: db,
		},
	}
}

func (uc *UseCase) Add(movieId uint, userId *uint, ip net.IP) (*View, error) {
	if userId == nil && ip == nil {
		return nil, errors.New("view: user ip and ip is nil")
	}

	lastViewMinTime := time.Now().Add(-24 * time.Hour)
	_, err := uc.ViewInteractor.GetLast(movieId, userId, ip.String(), lastViewMinTime)
	if err == nil {
		return nil, errors.New("view: too early to add a view")
	}

	view := &View{}
	view.UserID = userId
	view.IP = ip.String()

	return view, uc.ViewInteractor.Create(view)
}

func (uc *UseCase) GetCount(movieId uint) (int64, error) {
	return uc.ViewInteractor.GetCount(movieId)
}

func (uc *UseCase) GetMultipleCount(movieIds []uint) (map[uint]int64, error) {
	return uc.ViewInteractor.GetCountMultiple(movieIds)
}
