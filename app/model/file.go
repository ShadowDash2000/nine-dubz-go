package model

import "gorm.io/gorm"

type File struct {
	gorm.Model   `json:"-"`
	ID           uint
	Name         string `json:"name" gorm:"not null"`
	Extension    string `json:"extension" gorm:"not null"`
	OriginalName string `json:"originalName" gorm:"not null"`
	Path         string `json:"path" gorm:"not null"`
	Size         int64  `json:"size" gorm:"not null"`
}
