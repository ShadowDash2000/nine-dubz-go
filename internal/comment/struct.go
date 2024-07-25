package comment

import (
	"gorm.io/gorm"
	"nine-dubz/internal/movie"
	"nine-dubz/internal/user"
	"time"
)

type Comment struct {
	gorm.Model
	Text     string      `json:"text"`
	MovieID  uint        `json:"-"`
	Movie    movie.Movie `json:"-"`
	UserID   uint        `json:"-"`
	User     user.User   `json:"user"`
	ParentID *uint       `json:"-"`
	Parent   *Comment    `json:"parent,omitempty"`
}

type AddRequest struct {
	Text string `json:"text"`
}

type GetResponse struct {
	CreatedAt time.Time               `json:"createdAt"`
	Text      string                  `json:"text"`
	User      *user.GetPublicResponse `json:"user"`
	Parent    *GetResponse            `json:"parent,omitempty"`
}

func NewGetResponse(comment *Comment) *GetResponse {
	response := &GetResponse{
		CreatedAt: comment.CreatedAt,
		Text:      comment.Text,
		User:      user.NewGetPublicResponse(&comment.User),
	}
	if comment.Parent != nil {
		response.Parent = NewGetResponse(comment.Parent)
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
