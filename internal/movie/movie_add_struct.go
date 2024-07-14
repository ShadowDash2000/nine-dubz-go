package movie

import (
	"nine-dubz/internal/file"
)

type AddRequest struct {
	Poster      string     `json:"poster"`
	Name        string     `json:"name"`
	DubDate     string     `json:"dubDate"`
	VoiceActors string     `json:"voiceActors"`
	Genre       string     `json:"genre"`
	Video       *file.File `json:"video"`
}

func NewAddRequest(movie *AddRequest) *Movie {
	return &Movie{
		Poster:      movie.Poster,
		Name:        movie.Name,
		DubDate:     movie.DubDate,
		VoiceActors: movie.VoiceActors,
		Genre:       movie.Genre,
		Video:       movie.Video,
	}
}
