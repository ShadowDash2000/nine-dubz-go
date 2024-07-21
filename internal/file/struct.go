package file

import "gorm.io/gorm"

type File struct {
	gorm.Model   `json:"-"`
	ID           uint
	Name         string `json:"name" gorm:"not null"`
	Extension    string `json:"extension" gorm:"not null"`
	OriginalName string `json:"originalName" gorm:"not null"`
	Size         int64  `json:"size" gorm:"not null"`
	Type         string `json:"-" gorm:"not null"`
}

type UploadStatus struct {
	Status int    `json:"status"`
	Error  string `json:"error,omitempty"`
}
