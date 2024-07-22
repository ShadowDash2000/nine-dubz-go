package comment

import (
	"gorm.io/gorm"
	"nine-dubz/internal/movie"
	"nine-dubz/internal/user"
)

type Comment struct {
	gorm.Model
	Text      string      `json:"text"`
	UserId    uint        `json:"-"`
	User      user.User   `json:"user" gorm:"foreignKey:UserID,references:ID"`
	ParentId  *uint       `json:"-"`
	Parent    *Comment    `json:"parent" gorm:"foreignKey:ParentID,references:ID"`
	MovieCode string      `json:"-"`
	Movie     movie.Movie `json:"movie" gorm:"foreignKey:MovieCode,references:Code"`
}
