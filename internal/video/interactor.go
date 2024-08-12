package video

type Interactor interface {
	Create(video *Video) error
	GetWhere(where interface{}) (*Video, error)
	Delete(id uint) error
}
