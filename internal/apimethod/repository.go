package apimethod

import "gorm.io/gorm"

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) GetWhere(apiMethod *ApiMethod, where map[string]interface{}) error {
	result := r.DB.Where(where).First(apiMethod)

	return result.Error
}
