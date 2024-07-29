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

	if utf8.RuneCountInString(text) == 0 {
		return errors.New("comment text too short")
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
		parentComment, err := uc.CommentInteractor.Get(
			map[string]interface{}{"id": parentCommentId},
			"",
			"",
			&pagination.Pagination{
				Limit:  0,
				Offset: 0,
			},
		)
		if err != nil {
			return err
		} else if parentComment.Parent != nil {
			return errors.New("parent comment already exists")
		}

		comment.ParentID = &parentCommentId
	}

	return uc.CommentInteractor.Create(comment)
}

func (uc *UseCase) Get(userId uint, movieCode string, commentId uint, paginationSub *pagination.Pagination) (*GetResponse, error) {
	if paginationSub.Limit > 10 {
		paginationSub.Limit = 10
	}

	movieResponse, err := uc.MovieUseCase.Get(userId, movieCode)
	if err != nil {
		return nil, err
	}

	comment, err := uc.CommentInteractor.Get(
		map[string]interface{}{
			"id":        commentId,
			"movie_id":  movieResponse.ID,
			"parent_id": nil,
		},
		"created_at desc",
		"created_at asc",
		paginationSub,
	)
	if err != nil {
		return nil, err
	}

	return NewGetResponse(comment), nil
}

func (uc *UseCase) GetMultiple(userId uint, movieCode string, paginationMain *pagination.Pagination) (*[]GetResponse, error) {
	if paginationMain.Limit > 20 {
		paginationMain.Limit = 20
	}

	movieResponse, err := uc.MovieUseCase.Get(userId, movieCode)
	if err != nil {
		return nil, err
	}

	comments, err := uc.CommentInteractor.GetMultiple(
		map[string]interface{}{
			"movie_id":  movieResponse.ID,
			"parent_id": nil,
		},
		"created_at desc",
		"created_at asc",
		paginationMain,
		&pagination.Pagination{
			Limit:  10,
			Offset: 0,
		},
	)
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
