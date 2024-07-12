package payload

import "nine-dubz/app/model"

func NewMoviePayload(movie *model.Movie) *model.Movie {
	return &model.Movie{
		Poster:      movie.Poster,
		Name:        movie.Name,
		DubDate:     movie.DubDate,
		VoiceActors: movie.VoiceActors,
		Genre:       movie.Genre,
		Video:       movie.Video,
	}
}
