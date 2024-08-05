package view

import "time"

type Interactor interface {
	Create(view *View) error
	GetLast(movieId uint, userId *uint, ip string, time time.Time) (View, error)
	GetCount(movieId uint) (*int64, error)
}
