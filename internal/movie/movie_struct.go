package movie

import (
	"gorm.io/gorm"
	"nine-dubz/internal/file"
)

type Movie struct {
	*gorm.Model
	ID          uint       `json:"ID"`
	IsPublished bool       `json:"-" gorm:"default:false"`
	Poster      string     `json:"poster,omitempty"`
	Name        string     `json:"name"`
	DubDate     string     `json:"dubDate"`
	VoiceActors string     `json:"voiceActors"`
	Genre       string     `json:"genre"`
	VideoId     *uint      `json:"-" gorm:"foreignKey:VideoId;references:ID;OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Video       *file.File `json:"video"`
}
