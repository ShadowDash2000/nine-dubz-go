package movie

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"gorm.io/gorm"
	"io"
	"nine-dubz/internal/file"
	"nine-dubz/internal/pagination"
	"strings"
	"time"
)

type UseCase struct {
	MovieInteractor Interactor
	FileUseCase     *file.UseCase
}

func New(db *gorm.DB, fuc *file.UseCase) *UseCase {
	return &UseCase{
		MovieInteractor: &Repository{
			DB: db,
		},
		FileUseCase: fuc,
	}
}

func (uc *UseCase) Add(movieAddRequest *AddRequest) (*AddResponse, error) {
	movie := NewAddRequest(movieAddRequest)

	hasher := sha256.New()
	hasher.Write([]byte(time.Now().String()))
	hash := hasher.Sum(nil)

	encoded := base64.URLEncoding.EncodeToString(hash)
	encoded = strings.TrimRight(encoded, "=")

	movie.Code = encoded[:11]

	err := uc.MovieInteractor.Add(movie)
	if err != nil {
		return nil, err
	}

	return NewAddResponse(movie), nil
}

func (uc *UseCase) Delete(userId uint, code string) error {
	return uc.MovieInteractor.Delete(userId, code)
}

func (uc *UseCase) DeleteMultiple(userId uint, movies *[]DeleteRequest) error {
	for _, movie := range *movies {
		if err := uc.Delete(userId, movie.Code); err != nil {
			return err
		}
	}

	return nil
}

func (uc *UseCase) UpdateVideo(movie *VideoUpdateRequest) error {
	movieRequest := NewVideoUpdateRequest(movie)
	return uc.MovieInteractor.UpdatesWhere(movieRequest, map[string]interface{}{"code": movie.Code})
}

func (uc *UseCase) UpdateByUserId(userId uint, movie *UpdateRequest) error {
	movieRequest := NewUpdateRequest(movie)

	if movie.PreviewHeader.Size > 0 {
		buff := make([]byte, 512)
		_, err := movie.Preview.Read(buff)
		if err != nil {
			return err
		}
		_, err = movie.Preview.Seek(0, io.SeekStart)
		if err != nil {
			return err
		}
		isCorrectType, _ := uc.FileUseCase.VerifyFileType(buff, []string{"image/jpeg", "image/png", "image/webp"})
		if !isCorrectType {
			return errors.New("invalid file type")
		}

		preview, err := uc.FileUseCase.SaveFile("upload/video_previews/", movie.PreviewHeader.Filename, movie.Preview)
		if err != nil {
			return err
		}

		movieRequest.Preview = preview
	}

	return uc.MovieInteractor.UpdatesWhere(movieRequest, map[string]interface{}{"code": movie.Code, "user_id": userId})
}

func (uc *UseCase) UpdatePublishStatus(userId uint, movie *UpdatePublishStatusRequest) error {
	movieRequest := NewUpdatePublishStatusRequest(movie)
	return uc.MovieInteractor.UpdatesWhere(movieRequest, map[string]interface{}{"code": movie.Code, "user_id": userId})
}

func (uc *UseCase) Get(movieCode string) (*GetResponse, error) {
	movie, err := uc.MovieInteractor.Get(movieCode)
	if err != nil {
		return nil, err
	}

	return NewGetResponse(movie), nil
}

func (uc *UseCase) CheckByUser(userId uint, code string) bool {
	_, err := uc.MovieInteractor.GetWhere(code, map[string]interface{}{"user_id": userId})
	if err != nil {
		return false
	}

	return true
}

func (uc *UseCase) GetMultipleByUserId(userId uint, pagination *pagination.Pagination) ([]*GetForUserResponse, error) {
	if pagination.Limit > 20 {
		pagination.Limit = 20
	}

	movies, err := uc.MovieInteractor.GetMultipleByUserId(userId, pagination)
	if err != nil {
		return nil, err
	}

	if len(*movies) == 0 {
		return nil, err
	}

	var moviesPayload []*GetForUserResponse
	for _, movie := range *movies {
		moviesPayload = append(moviesPayload, NewGetForUserResponse(&movie))
	}

	return moviesPayload, nil
}

func (uc *UseCase) GetMultiple(pagination *pagination.Pagination) ([]*GetResponse, error) {
	if pagination.Limit > 20 {
		pagination.Limit = 20
	}

	movies, err := uc.MovieInteractor.GetMultiple(pagination)
	if err != nil {
		return nil, err
	}

	if len(*movies) == 0 {
		return nil, err
	}

	var moviesPayload []*GetResponse
	for _, movie := range *movies {
		moviesPayload = append(moviesPayload, NewGetResponse(&movie))
	}

	return moviesPayload, nil
}
