package movie

import (
	"nine-dubz/internal/file"
)

type UpdateRequest struct {
	ID          uint       `json:"id"`
	Poster      string     `json:"poster"`
	Name        string     `json:"name"`
	DubDate     string     `json:"dubDate"`
	VoiceActors string     `json:"voiceActors"`
	Genre       string     `json:"genre"`
	Video       *file.File `json:"video"`
}

type UpdateResponse struct {
	Poster      string     `json:"poster"`
	Name        string     `json:"name"`
	DubDate     string     `json:"dubDate"`
	VoiceActors string     `json:"voiceActors"`
	Genre       string     `json:"genre"`
	Video       *file.File `json:"video"`
}

func NewUpdateRequest(movie *UpdateRequest) *Movie {
	return &Movie{
		ID:          movie.ID,
		Poster:      movie.Poster,
		Name:        movie.Name,
		DubDate:     movie.DubDate,
		VoiceActors: movie.VoiceActors,
		Genre:       movie.Genre,
		Video:       movie.Video,
	}
}

func NewUpdateResponse(movie *Movie) *UpdateResponse {
	return &UpdateResponse{
		Poster:      movie.Poster,
		Name:        movie.Name,
		DubDate:     movie.DubDate,
		VoiceActors: movie.VoiceActors,
		Genre:       movie.Genre,
		Video:       movie.Video,
	}
}
