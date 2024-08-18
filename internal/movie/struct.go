package movie

import (
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"mime/multipart"
	"nine-dubz/internal/category"
	"nine-dubz/internal/file"
	"nine-dubz/internal/user"
	"nine-dubz/internal/video"
	"nine-dubz/internal/view"
	"time"
)

type Movie struct {
	gorm.Model
	ID                   uint              `json:"ID"`
	Status               string            `json:"-" gorm:"default:'uploading'"`
	CreatedAt            time.Time         `json:"createdAt"`
	Code                 string            `json:"code"`
	IsPublished          bool              `json:"-" gorm:"default:false"`
	Description          string            `json:"description"`
	PreviewId            *uint             `json:"-"`
	Preview              *file.File        `json:"preview,omitempty" gorm:"foreignKey:PreviewId;references:ID;"`
	PreviewWebpId        *uint             `json:"-"`
	PreviewWebp          *file.File        `json:"previewWebp,omitempty" gorm:"foreignKey:PreviewWebpId;references:ID;"`
	DefaultPreviewId     *uint             `json:"-"`
	DefaultPreview       *file.File        `json:"defaultPreview" gorm:"foreignKey:DefaultPreviewId;references:ID;"`
	DefaultPreviewWebpId *uint             `json:"-"`
	DefaultPreviewWebp   *file.File        `json:"defaultPreviewWebp" gorm:"foreignKey:DefaultPreviewWebpId;references:ID;"`
	Name                 string            `json:"name"`
	Videos               []video.Video     `gorm:"many2many:movie_videos"`
	UserId               uint              `json:"-"`
	User                 user.User         `json:"-" gorm:"foreignKey:UserId;references:ID"`
	Category             category.Category `gorm:"default:1;"`
	WebVttId             *uint             `json:"-"`
	WebVtt               *file.File        `json:"webVtt" gorm:"foreignKey:WebVttId;references:ID;"`
	Views                []view.View       `gorm:"-"`
}

const (
	StatusUploading = "uploading"
	StatusReady     = "ready"
)

type PoolItem struct {
	Ctx    context.Context
	Cancel context.CancelFunc
}

type VideoUploadHeader struct {
	Filename  string `json:"filename"`
	Size      int    `json:"size"`
	MovieCode string `json:"movieCode"`
	Token     string `json:"token"`
}

type AddRequest struct {
	UserId uint `json:"userId"`
}

type AddResponse struct {
	Code string `json:"code"`
}

func NewAddRequest(movie *AddRequest) *Movie {
	return &Movie{
		UserId: movie.UserId,
	}
}

func NewAddResponse(movie *Movie) *AddResponse {
	return &AddResponse{
		Code: movie.Code,
	}
}

type GetResponse struct {
	ID                 uint                    `json:"ID"`
	Code               string                  `json:"code"`
	CreatedAt          time.Time               `json:"createdAt"`
	Description        string                  `json:"description"`
	Preview            *file.File              `json:"preview"`
	PreviewWebp        *file.File              `json:"previewWebp"`
	DefaultPreview     *file.File              `json:"defaultPreview"`
	DefaultPreviewWebp *file.File              `json:"defaultPreviewWebp"`
	Name               string                  `json:"name"`
	Videos             []*video.GetResponse    `json:"videos"`
	Category           category.Category       `json:"category"`
	WebVtt             *file.File              `json:"webVtt"`
	User               *user.GetPublicResponse `json:"user"`
	Subscribed         bool                    `json:"subscribed"`
	Views              int64                   `json:"views"`
}

func NewGetResponse(movie *Movie) *GetResponse {
	return &GetResponse{
		ID:                 movie.ID,
		Code:               movie.Code,
		CreatedAt:          movie.CreatedAt,
		Description:        movie.Description,
		Preview:            movie.Preview,
		PreviewWebp:        movie.PreviewWebp,
		DefaultPreview:     movie.DefaultPreview,
		DefaultPreviewWebp: movie.DefaultPreviewWebp,
		Name:               movie.Name,
		Videos:             video.NewGetResponseMultiple(movie.Videos),
		Category:           movie.Category,
		WebVtt:             movie.WebVtt,
		User:               user.NewGetPublicResponse(&movie.User),
	}
}

type GetForUserResponse struct {
	IsPublished        bool                 `json:"isPublished"`
	Code               string               `json:"code"`
	CreatedAt          time.Time            `json:"createdAt"`
	Description        string               `json:"description"`
	Preview            *file.File           `json:"preview"`
	PreviewWebp        *file.File           `json:"previewWebp"`
	DefaultPreview     *file.File           `json:"defaultPreview"`
	DefaultPreviewWebp *file.File           `json:"defaultPreviewWebp"`
	Name               string               `json:"name"`
	Videos             []*video.GetResponse `json:"videos"`
}

func NewGetForUserResponse(movie *Movie) *GetForUserResponse {
	return &GetForUserResponse{
		IsPublished:        movie.IsPublished,
		Code:               movie.Code,
		CreatedAt:          movie.CreatedAt,
		Description:        movie.Description,
		Preview:            movie.Preview,
		PreviewWebp:        movie.PreviewWebp,
		DefaultPreview:     movie.DefaultPreview,
		DefaultPreviewWebp: movie.DefaultPreviewWebp,
		Name:               movie.Name,
		Videos:             video.NewGetResponseMultiple(movie.Videos),
	}
}

type VideoUpdateRequest struct {
	Status             string     `json:"-"`
	Name               string     `json:"name"`
	Code               string     `json:"code"`
	DefaultPreview     *file.File `json:"defaultPreview"`
	DefaultPreviewWebp *file.File `json:"defaultPreviewWebp"`
	WebVtt             *file.File `json:"webVtt"`
}

func NewVideoUpdateRequest(movie *VideoUpdateRequest) *Movie {
	return &Movie{
		Status:             movie.Status,
		Name:               movie.Name,
		Code:               movie.Code,
		DefaultPreview:     movie.DefaultPreview,
		DefaultPreviewWebp: movie.DefaultPreviewWebp,
		WebVtt:             movie.WebVtt,
	}
}

type UpdateRequest struct {
	Code          string                `json:"code"`
	IsPublished   bool                  `json:"isPublished,omitempty"`
	Description   string                `json:"description,omitempty"`
	Preview       multipart.File        `json:"preview,omitempty"`
	PreviewHeader *multipart.FileHeader `json:"-"`
	RemovePreview bool                  `json:"-"`
	Name          string                `json:"name,omitempty"`
	Category      category.Category     `json:"category,omitempty"`
}

func NewUpdateRequest(movie *UpdateRequest) *Movie {
	return &Movie{
		Code:        movie.Code,
		IsPublished: movie.IsPublished,
		Description: movie.Description,
		Name:        movie.Name,
		Category:    movie.Category,
	}
}

type DeleteRequest struct {
	Code string `json:"code"`
}

type UpdatePublishStatusRequest struct {
	Code        string `json:"code"`
	IsPublished bool   `json:"isPublished"`
}

func NewUpdatePublishStatusRequest(movie *UpdatePublishStatusRequest) *Movie {
	return &Movie{
		Code:        movie.Code,
		IsPublished: movie.IsPublished,
	}
}
