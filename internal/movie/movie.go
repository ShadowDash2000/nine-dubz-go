package movie

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/alitto/pond"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"io"
	"nine-dubz/internal/file"
	"nine-dubz/internal/pagination"
	"nine-dubz/pkg/ffmpegthumbs"
	"nine-dubz/pkg/webvtt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

type UseCase struct {
	MovieInteractor Interactor
	Pool            *pond.WorkerPool
	FileUseCase     *file.UseCase
}

func New(db *gorm.DB, pool *pond.WorkerPool, fuc *file.UseCase) *UseCase {
	return &UseCase{
		MovieInteractor: &Repository{
			DB: db,
		},
		Pool:        pool,
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

func (uc *UseCase) DeleteMovieFiles(movie *Movie) {
	if movie.Video != nil {
		uc.FileUseCase.RemoveFile(movie.Video.Name)
	}
	if movie.VideoShakal != nil {
		uc.FileUseCase.RemoveFile(movie.VideoShakal.Name)
	}
	if movie.Video360 != nil {
		uc.FileUseCase.RemoveFile(movie.Video360.Name)
	}
	if movie.Video480 != nil {
		uc.FileUseCase.RemoveFile(movie.Video480.Name)
	}
	if movie.Video720 != nil {
		uc.FileUseCase.RemoveFile(movie.Video720.Name)
	}
	if movie.WebVtt != nil {
		uc.FileUseCase.RemoveFile(movie.WebVtt.Name)
	}
	if len(movie.WebVttImages) > 0 {
		imagesNames := strings.Split(movie.WebVttImages, ";")
		for _, imageName := range imagesNames {
			uc.FileUseCase.RemoveFile(imageName)
		}
	}
}

func (uc *UseCase) Delete(userId uint, code string) error {
	movie, err := uc.MovieInteractor.GetWhere(code, map[string]interface{}{
		"user_id": userId,
	})
	if err != nil {
		return err
	}

	err = uc.MovieInteractor.Delete(userId, code)
	if err != nil {
		return err
	}

	go uc.DeleteMovieFiles(movie)

	return nil
}

func (uc *UseCase) DeleteMultiple(userId uint, movies *[]DeleteRequest) error {
	for _, movie := range *movies {
		if err := uc.Delete(userId, movie.Code); err != nil {
			return err
		}
	}

	return nil
}

func (uc *UseCase) SaveVideo(userId uint, header *VideoUploadHeader, conn *websocket.Conn) error {
	movieUpdateRequest := &VideoUpdateRequest{
		Code: header.MovieCode,
		Name: header.Filename,
	}
	rowsAffected, err := uc.UpdateVideo(movieUpdateRequest)
	if err != nil {
		return err
	} else if rowsAffected == 0 {
		return errors.New("failed to update video")
	}

	tmpFile, err := uc.FileUseCase.WriteFileFromSocket([]string{"video/mp4"}, header.Size, conn)
	if err != nil {
		uc.Delete(userId, header.MovieCode)
		return err
	}

	tmpFile, err = os.Open(tmpFile.Name())
	if err != nil {
		return errors.New("failed to open file tmp file")
	}

	tmpSavedFile, err := uc.FileUseCase.SaveFile(tmpFile, tmpFile.Name(), int64(header.Size), "private")
	if err != nil {
		return err
	}
	tmpFile.Close()

	movieUpdateRequest = &VideoUpdateRequest{
		Code:     header.MovieCode,
		VideoTmp: tmpSavedFile,
	}
	rowsAffected, err = uc.UpdateVideo(movieUpdateRequest)
	if err != nil || rowsAffected == 0 {
		uc.DeleteMovieFiles(NewVideoUpdateRequest(movieUpdateRequest))
		return errors.New("failed to update tmp video")
	}

	uc.Pool.Submit(func() {
		uc.PostProcessVideo(header, tmpFile)
	})

	return nil
}

func (uc *UseCase) RetryVideoPostProcess() {
	movies, err := uc.MovieInteractor.GetWhereMultiple(&pagination.Pagination{
		Limit:  -1,
		Offset: -1,
	}, map[string]interface{}{"status": "uploading"})
	if err != nil {
		return
	}

	for _, movie := range *movies {
		if movie.VideoTmp == nil {
			uc.Delete(movie.UserId, movie.Code)
			continue
		}
		tmpFile, err := os.Open(movie.VideoTmp.OriginalName)
		if err != nil {
			uc.Delete(movie.UserId, movie.Code)
			continue
		}

		header := &VideoUploadHeader{
			MovieCode: movie.Code,
		}

		uc.Pool.Submit(func() {
			uc.PostProcessVideo(header, tmpFile, movie)
		})
	}
}

func (uc *UseCase) PostProcessVideo(header *VideoUploadHeader, tmpFile *os.File, movies ...Movie) error {
	movie := &Movie{}
	if len(movies) == 1 {
		movie = &movies[0]
	}

	resizedVideoPath := filepath.Join("upload/resize", header.MovieCode)
	thumbsPath := filepath.Join("upload/thumbs/", header.MovieCode)

	if movie.Video == nil {
		bitrate, err := ffmpegthumbs.GetVideoBitrate(tmpFile.Name())
		if err != nil || bitrate <= 0 {
			bitrate = 30000
		}

		ffmpegthumbs.ToWebm(tmpFile.Name(), "31", "1", strconv.Itoa(bitrate), resizedVideoPath, "orig")
	}

	origWebmPath := filepath.Join(resizedVideoPath, "orig.webm")
	origWebm, err := os.Open(origWebmPath)
	if err != nil {
		return err
	}
	origWebmInfo, _ := os.Stat(origWebm.Name())

	defer os.RemoveAll(resizedVideoPath)
	defer os.RemoveAll(thumbsPath)
	defer os.Remove(tmpFile.Name())

	origWebmFile, err := uc.FileUseCase.SaveFile(origWebm, origWebmInfo.Name(), origWebmInfo.Size(), "private")
	if err != nil {
		return err
	}

	origWebm.Close()

	movieUpdateRequest := &VideoUpdateRequest{
		Code:  header.MovieCode,
		Video: origWebmFile,
	}
	rowsAffected, err := uc.UpdateVideo(movieUpdateRequest)
	if err != nil || rowsAffected == 0 {
		uc.DeleteMovieFiles(NewVideoUpdateRequest(movieUpdateRequest))
		return errors.New("failed to update video")
	}

	if movie.WebVtt == nil {
		// Thumbs
		uc.CreateThumbnails(thumbsPath, header, tmpFile)
	}

	// Resize
	uc.CreateResizedVideos(origWebmPath, resizedVideoPath, header, tmpFile, *movie)

	// Check if movie was deleted while post-processing
	movie, err = uc.MovieInteractor.Get(header.MovieCode)
	if err == nil {
		if movie.DeletedAt.Valid {
			uc.DeleteMovieFiles(movie)
			return errors.New("movie was deleted while post-processing")
		}

		if movie.VideoTmp != nil {
			uc.FileUseCase.RemoveFile(movie.VideoTmp.Name)
		}
	}

	return nil
}

func (uc *UseCase) CreateThumbnails(thumbsPath string, header *VideoUploadHeader, tmpFile *os.File) error {
	thumbsWebvttPath := "/api/file/"
	var imagesNames []string
	imagesFilePath := make([]string, 0)
	frameDuration := 10

	err := ffmpegthumbs.SplitVideoToThumbnails(tmpFile.Name(), thumbsPath, frameDuration)
	var defaultPreview *file.File
	if err != nil {
		fmt.Println(err)
	} else {
		items, _ := os.ReadDir(thumbsPath)
		defaultPreviewPos := len(items) / 2
		for i, item := range items {
			if item.IsDir() {
				continue
			}

			imageFile, _ := os.Open(filepath.Join(thumbsPath, item.Name()))
			imageFileInfo, _ := imageFile.Stat()
			savedImageFile, _ := uc.FileUseCase.SaveFile(imageFile, item.Name(), imageFileInfo.Size(), "public")
			imagesFilePath = append(imagesFilePath, filepath.Join(thumbsWebvttPath, savedImageFile.Name))
			imagesNames = append(imagesNames, savedImageFile.Name)

			if defaultPreviewPos == i+1 {
				defaultPreview = savedImageFile
			}

			imageFile.Close()
		}
	}

	var savedVttFile *file.File
	if len(imagesFilePath) > 0 {
		videoDuration, _ := ffmpegthumbs.GetVideoDuration(tmpFile.Name())
		vttFile, _ := webvtt.CreateFromFilePaths(imagesFilePath, thumbsPath, videoDuration, frameDuration)
		vttFile, _ = os.Open(vttFile.Name())
		savedVttFile, _ = uc.FileUseCase.SaveFile(vttFile, vttFile.Name(), 0, "public")
		vttFile.Close()
	}

	movieUpdateRequest := &VideoUpdateRequest{
		Code:           header.MovieCode,
		DefaultPreview: defaultPreview,
		WebVtt:         savedVttFile,
		WebVttImages:   strings.Join(imagesNames, ";"),
	}
	rowsAffected, err := uc.UpdateVideo(movieUpdateRequest)
	if err != nil {
		return err
	} else if rowsAffected == 0 {
		return errors.New("failed to update video")
	}

	return nil
}

func (uc *UseCase) CreateResizedVideos(origWebmPath, resizedVideoPath string, header *VideoUploadHeader, tmpFile *os.File, movies ...Movie) error {
	var videoShakal *file.File
	var video360 *file.File
	var video480 *file.File
	var video720 *file.File

	movie := &Movie{}
	if len(movies) == 1 {
		movie = &movies[0]
	}

	_, videoHeight, _ := ffmpegthumbs.GetVideoSize(origWebmPath)
	// 0P
	if videoHeight > 0 && movie.VideoShakal == nil {
		ffmpegthumbs.Resize(240, "50", "5", "5", tmpFile.Name(), resizedVideoPath, "shakal")
		resizedWebm, err := os.Open(filepath.Join(resizedVideoPath, "shakal.webm"))
		if err == nil {
			resizedWebmInfo, _ := os.Stat(resizedWebm.Name())
			resizedWebmFile, err := uc.FileUseCase.SaveFile(resizedWebm, resizedWebmInfo.Name(), resizedWebmInfo.Size(), "private")
			if err == nil {
				videoShakal = resizedWebmFile
			}
			resizedWebm.Close()
		}
	}
	// 360p
	if videoHeight > 360 && movie.Video360 == nil {
		ffmpegthumbs.Resize(360, "33", "3", "900k", tmpFile.Name(), resizedVideoPath, "360")
		resizedWebm, err := os.Open(filepath.Join(resizedVideoPath, "360.webm"))
		if err == nil {
			resizedWebmInfo, _ := os.Stat(resizedWebm.Name())
			resizedWebmFile, err := uc.FileUseCase.SaveFile(resizedWebm, resizedWebmInfo.Name(), resizedWebmInfo.Size(), "private")
			if err == nil {
				video360 = resizedWebmFile
			}
			resizedWebm.Close()
		}
	}
	// 480p
	if videoHeight > 480 && movie.Video480 == nil {
		ffmpegthumbs.Resize(480, "33", "3", "1000k", tmpFile.Name(), resizedVideoPath, "480")
		resizedWebm, err := os.Open(filepath.Join(resizedVideoPath, "480.webm"))
		if err == nil {
			resizedWebmInfo, _ := os.Stat(resizedWebm.Name())
			resizedWebmFile, err := uc.FileUseCase.SaveFile(resizedWebm, resizedWebmInfo.Name(), resizedWebmInfo.Size(), "private")
			if err == nil {
				video480 = resizedWebmFile
			}
			resizedWebm.Close()
		}
	}
	// 720p
	if videoHeight > 720 && movie.Video720 == nil {
		ffmpegthumbs.Resize(720, "32", "2", "1800k", tmpFile.Name(), resizedVideoPath, "720")
		resizedWebm, err := os.Open(filepath.Join(resizedVideoPath, "720.webm"))
		if err == nil {
			resizedWebmInfo, _ := os.Stat(resizedWebm.Name())
			resizedWebmFile, err := uc.FileUseCase.SaveFile(resizedWebm, resizedWebmInfo.Name(), resizedWebmInfo.Size(), "private")
			if err == nil {
				video720 = resizedWebmFile
			}
			resizedWebm.Close()
		}
	}

	movieUpdateRequest := &VideoUpdateRequest{
		Status:      "ready",
		Code:        header.MovieCode,
		Video360:    video360,
		Video480:    video480,
		Video720:    video720,
		VideoShakal: videoShakal,
	}
	rowsAffected, err := uc.UpdateVideo(movieUpdateRequest)
	if err != nil {
		return err
	} else if rowsAffected == 0 {
		return errors.New("failed to update video")
	}

	return nil
}

func (uc *UseCase) UpdateVideo(movie *VideoUpdateRequest) (int64, error) {
	movieRequest := NewVideoUpdateRequest(movie)
	return uc.MovieInteractor.UpdatesWhere(movieRequest, map[string]interface{}{"code": movie.Code})
}

func (uc *UseCase) UpdateByUserId(userId uint, movie *UpdateRequest) error {
	movieRequest := NewUpdateRequest(movie)

	if utf8.RuneCountInString(movie.Name) > 130 {
		return errors.New("movie name too long")
	}
	if utf8.RuneCountInString(movie.Description) > 5000 {
		return errors.New("movie description too long")
	}

	if movie.PreviewHeader != nil && movie.PreviewHeader.Size > 0 {
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

		preview, err := uc.FileUseCase.SaveFile(movie.Preview, movie.PreviewHeader.Filename, movie.PreviewHeader.Size, "public")
		if err != nil {
			return err
		}

		movieRequest.Preview = preview
	}

	rowsAffected, err := uc.MovieInteractor.UpdatesWhere(movieRequest, map[string]interface{}{"code": movie.Code, "user_id": userId})
	if err != nil {
		return err
	} else if rowsAffected == 0 {
		return errors.New("movie not found")
	}

	return nil
}

func (uc *UseCase) UpdatePublishStatus(userId uint, movie *UpdatePublishStatusRequest) (int64, error) {
	movieRequest := NewUpdatePublishStatusRequest(movie)
	return uc.MovieInteractor.UpdatesWhere(movieRequest, map[string]interface{}{"code": movie.Code, "user_id": userId})
}

func (uc *UseCase) Get(userId uint, code string) (*GetResponse, error) {
	movie, err := uc.MovieInteractor.Get(code)
	if err != nil {
		return nil, err
	}

	if movie.IsPublished {
		return NewGetResponse(movie), nil
	}

	if movie.UserId == userId {
		return NewGetResponse(movie), nil
	}

	return nil, errors.New("not allowed")
}

func (uc *UseCase) CheckMovieAccess(userId uint, code string) (*GetResponse, error) {
	movie, err := uc.MovieInteractor.Get(code)
	if err != nil {
		return nil, err
	}

	if movie.Video == nil {
		return nil, errors.New("no video")
	}

	if movie.IsPublished {
		return NewGetResponse(movie), nil
	}

	if movie.UserId == userId {
		return NewGetResponse(movie), nil
	}

	return nil, errors.New("not allowed")
}

func (uc *UseCase) CheckByUser(userId uint, code string) bool {
	_, err := uc.MovieInteractor.GetWhere(code, map[string]interface{}{"user_id": userId, "video_id": nil})
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

func (uc *UseCase) GetForUser(userId uint, code string) (*GetForUserResponse, error) {
	movie, err := uc.MovieInteractor.GetWhere(code, map[string]interface{}{"user_id": userId})
	if err != nil {
		return nil, err
	}

	return NewGetForUserResponse(movie), nil
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
