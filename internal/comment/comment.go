package comment

import (
	"errors"
	"gorm.io/gorm"
	"nine-dubz/internal/movie"
	"nine-dubz/internal/pagination"
	"unicode/utf8"
)

type UseCase struct {
	CommentInteractor Interactor
	MovieUseCase      *movie.UseCase
}

func New(db *gorm.DB, muc *movie.UseCase) *UseCase {
	return &UseCase{
		CommentInteractor: &Repository{
			DB: db,
		},
		MovieUseCase: muc,
	}
}

func (uc *UseCase) Add(userId uint, movieCode, text string, options ...uint) error {
	var parentCommentId uint
	if len(options) > 0 {
		parentCommentId = options[0]
	}

	if utf8.RuneCountInString(text) > 5000 {
		return errors.New("comment text too long")
	}

	movieResponse, err := uc.MovieUseCase.Get(userId, movieCode)
	if err != nil {
		return err
	}

	comment := &Comment{
		Text:    text,
		MovieID: movieResponse.ID,
		UserID:  userId,
	}

	if parentCommentId > 0 {
		parentComment, err := uc.Get(parentCommentId)
		if err != nil {
			return err
		} else if parentComment.Parent != nil {
			return errors.New("parent comment already exists")
		}

		comment.ParentID = &parentCommentId
	}

	return uc.CommentInteractor.Create(comment)
}

func (uc *UseCase) Get(commentId uint) (*GetResponse, error) {
	comment, err := uc.CommentInteractor.Get(commentId)
	if err != nil {
		return nil, err
	}

	return NewGetResponse(comment), nil
}

func (uc *UseCase) GetMultiple(userId uint, movieCode string, pagination *pagination.Pagination) (*[]GetResponse, error) {
	if pagination.Limit > 20 {
		pagination.Limit = 20
	}

	movieResponse, err := uc.MovieUseCase.Get(userId, movieCode)
	if err != nil {
		return nil, err
	}

	comments, err := uc.CommentInteractor.GetMultiple(map[string]interface{}{
		"movie_id": movieResponse.ID,
	}, pagination)
	if err != nil {
		return nil, err
	}

	if len(*comments) == 0 {
		return nil, err
	}

	return NewGetMultipleResponse(comments), nil
}

func (uc *UseCase) Delete(commentId, userId uint) (int64, error) {
	return uc.CommentInteractor.Delete(commentId, userId)
}
