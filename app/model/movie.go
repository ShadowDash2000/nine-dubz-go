package model

import "gorm.io/gorm"

type Movie struct {
	gorm.Model
	Poster      string `json:"poster"`
	Name        string `json:"name"`
	DubDate     string `json:"dubDate"`
	VoiceActors string `json:"voiceActors"`
	Genre       string `json:"genre"`
	VideoUrl    string `json:"videoUrl"`
}
