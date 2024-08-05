package view

import (
	"gorm.io/gorm"
	"time"
)

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) Create(view *View) error {
	return r.DB.Create(&view).Error
}

func (r *Repository) GetLast(movieId uint, userId *uint, ip string, time time.Time) (View, error) {
	view := View{}
	result := r.DB.Where(
		"movie_id = ? AND created_at > ? AND (user_id = ? OR ip = ?)",
		movieId, time, userId, ip,
	).First(&view)

	return view, result.Error
}

func (r *Repository) GetCount(movieId uint) (*int64, error) {
	var count int64
	result := r.DB.Model(&View{}).Where("movie_id = ?", movieId).Count(&count)

	return &count, result.Error
}
