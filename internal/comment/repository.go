package comment

import "gorm.io/gorm"

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) Create() (*Comment, error) {
	comment := &Comment{}
	result := r.DB.Create(&comment)

	return comment, result.Error
}
