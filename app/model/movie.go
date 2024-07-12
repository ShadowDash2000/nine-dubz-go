package model

import "gorm.io/gorm"

type Movie struct {
	*gorm.Model
	IsPublished bool   `json:"-" gorm:"default:false"`
	Poster      string `json:"poster,omitempty"`
	Name        string `json:"name"`
	DubDate     string `json:"dubDate"`
	VoiceActors string `json:"voiceActors"`
	Genre       string `json:"genre"`
	VideoId     *uint  `json:"-" gorm:"foreignKey:VideoId;references:ID;OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Video       *File  `json:"video"`
}
