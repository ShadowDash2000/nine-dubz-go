package movie

import (
	"gorm.io/gorm"
	"mime/multipart"
	"nine-dubz/internal/file"
	"nine-dubz/internal/user"
	"time"
)

type Movie struct {
	gorm.Model
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
	VideoShakalId    *uint      `json:"-"`
	VideoShakal      *file.File `json:"videoShakal" gorm:"foreignKey:VideoShakalId;references:ID;"`
	Video360Id       *uint      `json:"-"`
	Video360         *file.File `json:"video360" gorm:"foreignKey:Video360Id;references:ID;"`
	Video480Id       *uint      `json:"-"`
	Video480         *file.File `json:"video480" gorm:"foreignKey:Video480Id;references:ID;"`
	Video720Id       *uint      `json:"-"`
	Video720         *file.File `json:"video720" gorm:"foreignKey:Video720Id;references:ID;"`
	UserId           uint       `json:"-"`
	User             user.User  `json:"-" gorm:"foreignKey:UserId;references:ID"`
	WebVttId         *uint      `json:"-"`
	WebVtt           *file.File `json:"webVtt" gorm:"foreignKey:WebVttId;references:ID;"`
	WebVttImages     string     `json:"-"`
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
	IsPublished    bool       `json:"isPublished"`
	Code           string     `json:"code"`
	CreatedAt      time.Time  `json:"createdAt"`
	Description    string     `json:"description"`
	Preview        *file.File `json:"preview"`
	DefaultPreview *file.File `json:"defaultPreview"`
	Name           string     `json:"name"`
	Video          *file.File `json:"video"`
	VideoShakal    *file.File `json:"videoShakal"`
	Video360       *file.File `json:"video360"`
	Video480       *file.File `json:"video480"`
	Video720       *file.File `json:"video720"`
	WebVtt         *file.File `json:"webVtt"`
}

func NewGetResponse(movie *Movie) *GetResponse {
	return &GetResponse{
		IsPublished:    movie.IsPublished,
		Code:           movie.Code,
		CreatedAt:      movie.CreatedAt,
		Description:    movie.Description,
		Preview:        movie.Preview,
		DefaultPreview: movie.DefaultPreview,
		Name:           movie.Name,
		Video:          movie.Video,
		VideoShakal:    movie.VideoShakal,
		Video360:       movie.Video360,
		Video480:       movie.Video480,
		Video720:       movie.Video720,
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
	VideoShakal    *file.File `json:"videoShakal"`
	Video360       *file.File `json:"video360"`
	Video480       *file.File `json:"video480"`
	Video720       *file.File `json:"video720"`
	DefaultPreview *file.File `json:"defaultPreview"`
	WebVtt         *file.File `json:"webVtt"`
	WebVttImages   string     `json:"-"`
}

func NewVideoUpdateRequest(movie *VideoUpdateRequest) *Movie {
	return &Movie{
		Name:           movie.Name,
		Code:           movie.Code,
		Video:          movie.Video,
		VideoShakal:    movie.VideoShakal,
		Video360:       movie.Video360,
		Video480:       movie.Video480,
		Video720:       movie.Video720,
		DefaultPreview: movie.DefaultPreview,
		WebVtt:         movie.WebVtt,
		WebVttImages:   movie.WebVttImages,
	}
}

type UpdateRequest struct {
	Code          string                `json:"code"`
	IsPublished   bool                  `json:"isPublished,omitempty"`
	Description   string                `json:"description,omitempty"`
	Preview       multipart.File        `json:"preview,omitempty"`
	PreviewHeader *multipart.FileHeader `json:"-"`
	Name          string                `json:"name,omitempty"`
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
