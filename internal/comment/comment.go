package comment

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"nine-dubz/internal/movie"
	"nine-dubz/internal/pagination"
	"nine-dubz/internal/sort"
	"nine-dubz/internal/user"
	"regexp"
	"slices"
	"strconv"
	"unicode/utf8"
)

type UseCase struct {
	CommentInteractor Interactor
	MovieUseCase      *movie.UseCase
	UserUseCase       *user.UseCase
}

func New(db *gorm.DB, muc *movie.UseCase, uuc *user.UseCase) *UseCase {
	return &UseCase{
		CommentInteractor: &Repository{
			DB: db,
		},
		MovieUseCase: muc,
		UserUseCase:  uuc,
	}
}

func (uc *UseCase) Add(userId uint, movieCode, text string, options ...uint) error {
	var parentCommentId uint
	if len(options) > 0 {
		parentCommentId = options[0]
	}

	if utf8.RuneCountInString(text) == 0 {
		return errors.New("comment text is required")
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
		parentComments, err := uc.CommentInteractor.GetDistinctMultiple(
			map[string]interface{}{"id": parentCommentId},
			[]string{"parent_id"},
		)
		if err != nil {
			return err
		} else if len(parentComments) > 0 {
			if parentComments[0].ParentID != nil {
				return errors.New("parent comment already exists")
			}
		}

		comment.ParentID = &parentCommentId
	}

	return uc.CommentInteractor.Create(comment)
}

func (uc *UseCase) GetMultipleSubComments(userId uint, movieCode string, parentId uint, paginationMain *pagination.Pagination) (*[]GetSubCommentResponse, error) {
	if paginationMain.Limit > 10 {
		paginationMain.Limit = 10
	}

	movieResponse, err := uc.MovieUseCase.Get(userId, movieCode)
	if err != nil {
		return nil, err
	}

	comments, err := uc.CommentInteractor.GetMultiple(
		map[string]interface{}{
			"movie_id":  movieResponse.ID,
			"parent_id": parentId,
		},
		"created_at asc",
		"",
		paginationMain,
		&pagination.Pagination{
			Limit:  -1,
			Offset: -1,
		},
	)
	if err != nil {
		return nil, err
	}

	err = uc.Format(&comments)
	if err != nil {
		return nil, errors.New("comment: error while formatting comments")
	}

	return NewGetMultipleSubCommentResponse(&comments), nil
}

func (uc *UseCase) GetMultiple(userId uint, movieCode string, paginationMain *pagination.Pagination, sort *sort.Sort) (*GetMultipleResponse, error) {
	if paginationMain.Limit > 20 {
		paginationMain.Limit = 20
	}

	if !slices.Contains([]string{"created_at"}, sort.SortBy) {
		sort.SortBy = "created_at"
		sort.SortVal = "desc"
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
		fmt.Sprintf("%s %s", sort.SortBy, sort.SortVal),
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

	if len(comments) == 0 {
		return nil, err
	}

	commentsCount, err := uc.CommentInteractor.Count(map[string]interface{}{
		"movie_id":  movieResponse.ID,
		"parent_id": nil,
	})
	if err != nil {
		return nil, err
	}

	var commentsIds []uint
	for _, comment := range comments {
		commentsIds = append(commentsIds, comment.ID)
	}

	subComments, err := uc.CommentInteractor.GetDistinctMultiple(
		map[string]interface{}{"parent_id": commentsIds},
		[]string{"id", "parent_id"},
	)
	if err != nil {
		return nil, err
	}

	subCommentsCount := make(map[uint]int64)
	for _, subComment := range subComments {
		subCommentsCount[*subComment.ParentID] = subCommentsCount[*subComment.ParentID] + 1
	}

	if len(subCommentsCount) > 0 {
		for key, comment := range comments {
			if _, ok := subCommentsCount[comment.ID]; ok {
				comments[key].SubCommentsCount = subCommentsCount[comment.ID]
			}
		}
	}

	err = uc.Format(&comments)
	if err != nil {
		return nil, errors.New("comment: error while formatting comments")
	}

	return NewGetMultipleResponse(&comments, commentsCount), nil
}

func (uc *UseCase) Format(comments *[]Comment) error {
	r := regexp.MustCompile(`<@id:(\d*)>`)
	var userIds []uint
	type Text struct {
		UserIds  []uint
		Mentions *[]Mention
	}
	var texts []Text

	for i, comment := range *comments {
		commentText := Text{
			Mentions: &(*comments)[i].Mentions,
		}
		matches := r.FindAllStringSubmatch(comment.Text, -1)
		for _, match := range matches {
			userId, err := strconv.ParseUint(match[1], 10, 32)
			if err != nil {
				continue
			}
			if !slices.Contains(userIds, uint(userId)) {
				userIds = append(userIds, uint(userId))
			}
			if !slices.Contains(commentText.UserIds, uint(userId)) {
				commentText.UserIds = append(commentText.UserIds, uint(userId))
			}
		}

		for j, subComment := range comment.SubComments {
			subCommentsText := Text{
				Mentions: &(*comments)[i].SubComments[j].Mentions,
			}
			matches = r.FindAllStringSubmatch(subComment.Text, -1)
			for _, match := range matches {
				userId, err := strconv.ParseUint(match[1], 10, 32)
				if err != nil {
					continue
				}
				if !slices.Contains(userIds, uint(userId)) {
					userIds = append(userIds, uint(userId))
				}
				if !slices.Contains(subCommentsText.UserIds, uint(userId)) {
					subCommentsText.UserIds = append(subCommentsText.UserIds, uint(userId))
				}
			}

			texts = append(texts, subCommentsText)
		}

		texts = append(texts, commentText)
	}

	if len(userIds) == 0 {
		return nil
	}

	users, err := uc.UserUseCase.GetMultiple(
		map[string]interface{}{"id": userIds},
		[]string{"id", "name"},
	)
	if err != nil {
		return err
	}

	userNames := make(map[uint]string)
	for _, user := range users {
		userNames[user.ID] = user.Name
	}
	for _, text := range texts {
		i := 0
		for _, userId := range text.UserIds {
			if _, ok := userNames[userId]; ok {
				mention := fmt.Sprintf("@%s", userNames[userId])
				*text.Mentions = append(*text.Mentions, Mention{
					UserID:  userId,
					Mention: mention,
				})
				i = i + 1
			}
		}
	}

	return nil
}

func (uc *UseCase) Delete(commentId, userId uint) (int64, error) {
	return uc.CommentInteractor.Delete(commentId, userId)
}
