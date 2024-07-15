package movie

import (
	"gorm.io/gorm"
	"nine-dubz/internal/file"
	"nine-dubz/internal/user"
)

type Movie struct {
	*gorm.Model
	ID          uint       `json:"ID"`
	Code        string     `json:"code"`
	IsPublished bool       `json:"-" gorm:"default:false"`
	Poster      string     `json:"poster,omitempty"`
	Name        string     `json:"name"`
	DubDate     string     `json:"dubDate"`
	VoiceActors string     `json:"voiceActors"`
	Genre       string     `json:"genre"`
	VideoId     *uint      `json:"-" gorm:"foreignKey:VideoId;references:ID;OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Video       *file.File `json:"video"`
	UserId      uint       `json:"-" gorm:"foreignKey:UserID;references:ID"`
	User        user.User  `json:"-"`
}

type VideoUploadHeader struct {
	Filename  string `json:"filename"`
	Size      int    `json:"size"`
	MovieCode string `json:"movieCode"`
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
	Code        string     `json:"code"`
	Poster      string     `json:"poster"`
	Name        string     `json:"name"`
	DubDate     string     `json:"dubDate"`
	VoiceActors string     `json:"voiceActors"`
	Genre       string     `json:"genre"`
	Video       *file.File `json:"video"`
}

func NewGetResponse(movie *Movie) *GetResponse {
	return &GetResponse{
		Poster:      movie.Poster,
		Name:        movie.Name,
		DubDate:     movie.DubDate,
		VoiceActors: movie.VoiceActors,
		Genre:       movie.Genre,
		Video:       movie.Video,
	}
}

type GetForUserResponse struct {
	IsPublished bool       `json:"isPublished"`
	Code        string     `json:"code"`
	Poster      string     `json:"poster"`
	Name        string     `json:"name"`
	DubDate     string     `json:"dubDate"`
	VoiceActors string     `json:"voiceActors"`
	Genre       string     `json:"genre"`
	Video       *file.File `json:"video"`
}

func NewGetForUserResponse(movie *Movie) *GetForUserResponse {
	return &GetForUserResponse{
		IsPublished: movie.IsPublished,
		Poster:      movie.Poster,
		Name:        movie.Name,
		DubDate:     movie.DubDate,
		VoiceActors: movie.VoiceActors,
		Genre:       movie.Genre,
		Video:       movie.Video,
	}
}

type UpdateRequest struct {
	Code        string     `json:"code"`
	IsPublished bool       `json:"isPublished"`
	Poster      string     `json:"poster"`
	Name        string     `json:"name"`
	DubDate     string     `json:"dubDate"`
	VoiceActors string     `json:"voiceActors"`
	Genre       string     `json:"genre"`
	Video       *file.File `json:"video"`
}

func NewUpdateRequest(movie *UpdateRequest) *Movie {
	return &Movie{
		Code:        movie.Code,
		Poster:      movie.Poster,
		Name:        movie.Name,
		DubDate:     movie.DubDate,
		VoiceActors: movie.VoiceActors,
		Genre:       movie.Genre,
		Video:       movie.Video,
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
