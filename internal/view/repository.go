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

func (r *Repository) GetCount(movieId uint) (int64, error) {
	var count int64
	result := r.DB.Model(&View{}).Where("movie_id = ?", movieId).Count(&count)

	return count, result.Error
}

func (r *Repository) GetCountMultiple(movieIds []uint) (map[uint]int64, error) {
	var views []View
	result := r.DB.Select("movie_id").Where("movie_id IN ?", movieIds).Find(&views)

	counts := make(map[uint]int64)
	for _, view := range views {
		if _, ok := counts[view.MovieID]; ok {
			counts[view.MovieID] = counts[view.MovieID] + 1
		} else {
			counts[view.MovieID] = 1
		}
	}

	return counts, result.Error
}
