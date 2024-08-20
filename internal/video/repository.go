package video

import "gorm.io/gorm"

type Repository struct {
	DB *gorm.DB
}

func (r *Repository) Create(video *Video) error {
	return r.DB.Create(video).Error
}

func (r *Repository) GetWhere(where interface{}) (*Video, error) {
	video := &Video{}
	result := r.DB.Where(where).First(video)

	return video, result.Error
}

func (r *Repository) Delete(id uint) error {
	return r.DB.Delete(&Video{ID: id}).Error
}
