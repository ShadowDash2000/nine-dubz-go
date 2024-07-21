package movie

import (
	"gorm.io/gorm"
	"mime/multipart"
	"nine-dubz/internal/file"
	"nine-dubz/internal/user"
	"time"
)

type Movie struct {
	*gorm.Model
	ID               uint       `json:"ID"`
	CreatedAt        time.Time  `json:"createdAt"`
	Code             string     `json:"code"`
	IsPublished      bool       `json:"-" gorm:"default:false"`
	Description      string     `json:"description"`
	PreviewId        *uint      `json:"-"`
	Preview          *file.File `json:"preview,omitempty" gorm:"foreignKey:PreviewId;references:ID;"`
	DefaultPreviewId *uint      `json:"-"`
	DefaultPreview   *file.File `json:"defaultPreview" gorm:"foreignKey:DefaultPreviewId;references:ID;"`
	Name             string     `json:"name"`
	VideoId          *uint      `json:"-"`
	Video            *file.File `json:"video" gorm:"foreignKey:VideoId;references:ID;"`
	UserId           uint       `json:"-"`
	User             user.User  `json:"-" gorm:"foreignKey:UserId;references:ID"`
	WebVttId         *uint      `json:"-"`
	WebVtt           *file.File `json:"webVtt" gorm:"foreignKey:WebVttId;references:ID;"`
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
	Code           string     `json:"code"`
	CreatedAt      time.Time  `json:"createdAt"`
	Description    string     `json:"description"`
	Preview        *file.File `json:"preview"`
	DefaultPreview *file.File `json:"defaultPreview"`
	Name           string     `json:"name"`
	Video          *file.File `json:"video"`
	WebVtt         *file.File `json:"webVtt"`
}

func NewGetResponse(movie *Movie) *GetResponse {
	return &GetResponse{
		Code:           movie.Code,
		CreatedAt:      movie.CreatedAt,
		Description:    movie.Description,
		Preview:        movie.Preview,
		DefaultPreview: movie.DefaultPreview,
		Name:           movie.Name,
		Video:          movie.Video,
		WebVtt:         movie.WebVtt,
	}
}

type GetForUserResponse struct {
	IsPublished    bool       `json:"isPublished"`
	Code           string     `json:"code"`
	CreatedAt      time.Time  `json:"createdAt"`
	Description    string     `json:"description"`
	Preview        *file.File `json:"preview"`
	DefaultPreview *file.File `json:"defaultPreview"`
	Name           string     `json:"name"`
	Video          *file.File `json:"video"`
}

func NewGetForUserResponse(movie *Movie) *GetForUserResponse {
	return &GetForUserResponse{
		IsPublished:    movie.IsPublished,
		Code:           movie.Code,
		CreatedAt:      movie.CreatedAt,
		Description:    movie.Description,
		Preview:        movie.Preview,
		DefaultPreview: movie.DefaultPreview,
		Name:           movie.Name,
		Video:          movie.Video,
	}
}

type VideoUpdateRequest struct {
	Name           string     `json:"name"`
	Code           string     `json:"code"`
	Video          *file.File `json:"video"`
	DefaultPreview *file.File `json:"defaultPreview"`
	WebVtt         *file.File `json:"webVtt"`
}

func NewVideoUpdateRequest(movie *VideoUpdateRequest) *Movie {
	return &Movie{
		Name:           movie.Name,
		Code:           movie.Code,
		Video:          movie.Video,
		DefaultPreview: movie.DefaultPreview,
		WebVtt:         movie.WebVtt,
	}
}

type UpdateRequest struct {
	Code          string                `json:"code"`
	IsPublished   bool                  `json:"isPublished"`
	Description   string                `json:"description"`
	Preview       multipart.File        `json:"preview,omitempty"`
	PreviewHeader *multipart.FileHeader `json:"-"`
	Name          string                `json:"name"`
}

func NewUpdateRequest(movie *UpdateRequest) *Movie {
	return &Movie{
		Code:        movie.Code,
		IsPublished: movie.IsPublished,
		Description: movie.Description,
		Name:        movie.Name,
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
