package comment

import (
	"gorm.io/gorm"
	"nine-dubz/internal/movie"
	"nine-dubz/internal/user"
	"time"
)

type Comment struct {
	gorm.Model
	Text             string
	Mentions         []Mention `gorm:"-"`
	MovieID          uint      `gorm:"not null"`
	Movie            movie.Movie
	UserID           uint `gorm:"not null"`
	User             user.User
	ParentID         *uint
	Parent           *Comment
	SubCommentsCount int64     `gorm:"-"`
	SubComments      []Comment `gorm:"foreignKey:ParentID"`
}

type Mention struct {
	UserID  uint   `json:"userId"`
	Mention string `json:"mention"`
}

type AddRequest struct {
	Text string `json:"text"`
}

type GetResponse struct {
	ID               uint                    `json:"id"`
	CreatedAt        time.Time               `json:"createdAt"`
	Text             string                  `json:"text"`
	Mentions         []Mention               `json:"mentions,omitempty"`
	User             *user.GetPublicResponse `json:"user"`
	Parent           *GetResponse            `json:"-"`
	SubCommentsCount int64                   `json:"subCommentsCount"`
	SubComments      []GetSubCommentResponse `json:"subComments,omitempty"`
}

func NewGetResponse(comment *Comment) *GetResponse {
	response := &GetResponse{
		ID:               comment.ID,
		CreatedAt:        comment.CreatedAt,
		Text:             comment.Text,
		Mentions:         comment.Mentions,
		User:             user.NewGetPublicResponse(&comment.User),
		SubCommentsCount: comment.SubCommentsCount,
		SubComments:      *NewGetMultipleSubCommentResponse(&comment.SubComments),
	}

	return response
}

type GetSubCommentResponse struct {
	ID        uint                    `json:"id"`
	CreatedAt time.Time               `json:"createdAt"`
	Text      string                  `json:"text"`
	Mentions  []Mention               `json:"mentions,omitempty"`
	User      *user.GetPublicResponse `json:"user"`
}

func NewGetSubCommentResponse(comment *Comment) *GetSubCommentResponse {
	response := &GetSubCommentResponse{
		ID:        comment.ID,
		CreatedAt: comment.CreatedAt,
		Text:      comment.Text,
		Mentions:  comment.Mentions,
		User:      user.NewGetPublicResponse(&comment.User),
	}

	return response
}

func NewGetMultipleSubCommentResponse(comments *[]Comment) *[]GetSubCommentResponse {
	var commentsResponse []GetSubCommentResponse
	for _, comment := range *comments {
		commentsResponse = append(commentsResponse, *NewGetSubCommentResponse(&comment))
	}

	return &commentsResponse
}

type GetMultipleResponse struct {
	CommentsCount int64         `json:"commentsCount"`
	Comments      []GetResponse `json:"comments"`
}

func NewGetMultipleResponse(comments *[]Comment, commentsCount int64) *GetMultipleResponse {
	var commentsResponse GetMultipleResponse
	for _, comment := range *comments {
		commentsResponse.Comments = append(commentsResponse.Comments, *NewGetResponse(&comment))
	}

	if commentsCount > 0 {
		commentsResponse.CommentsCount = commentsCount
	}

	return &commentsResponse
}
