package comment

import (
	"gorm.io/gorm"
	"nine-dubz/internal/movie"
	"nine-dubz/internal/user"
	"time"
)

type Comment struct {
	gorm.Model
	Text        string
	MovieID     uint `gorm:"not null"`
	Movie       movie.Movie
	UserID      uint `gorm:"not null"`
	User        user.User
	ParentID    *uint
	Parent      *Comment
	SubComments []Comment `gorm:"foreignKey:ParentID"`
}

type AddRequest struct {
	Text string `json:"text"`
}

type GetResponse struct {
	CreatedAt   time.Time               `json:"createdAt"`
	Text        string                  `json:"text"`
	User        *user.GetPublicResponse `json:"user"`
	Parent      *GetResponse            `json:"-"`
	SubComments []GetResponse           `json:"subComments,omitempty"`
}

func NewGetResponse(comment *Comment) *GetResponse {
	response := &GetResponse{
		CreatedAt:   comment.CreatedAt,
		Text:        comment.Text,
		User:        user.NewGetPublicResponse(&comment.User),
		SubComments: *NewGetMultipleResponse(&comment.SubComments),
	}

	return response
}

func NewGetMultipleResponse(comments *[]Comment) *[]GetResponse {
	var commentsResponse []GetResponse
	for _, comment := range *comments {
		commentsResponse = append(commentsResponse, *NewGetResponse(&comment))
	}

	return &commentsResponse
}
