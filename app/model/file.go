package model

import "gorm.io/gorm"

type File struct {
	gorm.Model   `json:"-"`
	ID           uint
	Name         int64  `json:"name" gorm:"not null"`
	OriginalName string `json:"originalName" gorm:"not null"`
	Path         string `json:"path" gorm:"not null"`
}
