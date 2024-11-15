package movie

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"nine-dubz/internal/file"
	"nine-dubz/internal/pagination"
	"nine-dubz/internal/sorting"
	"nine-dubz/internal/subscription"
	"nine-dubz/internal/video"
	"nine-dubz/internal/view"
	"nine-dubz/pkg/ffmpegthumbs"
	"nine-dubz/pkg/language"
	"nine-dubz/pkg/webvtt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/alitto/pond"
	"github.com/aws/smithy-go/ptr"
	"github.com/gorilla/websocket"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type UseCase struct {
	MovieInteractor     Interactor
	SiteUrl             string
	Pool                *pond.WorkerPool
	VideoUseCase        *video.UseCase
	FileUseCase         *file.UseCase
	ViewUseCase         *view.UseCase
	SubscriptionUseCase *subscription.UseCase
	MoviePool           map[string]PoolItem
	Mutex               *sync.RWMutex
}

func New(db *gorm.DB, pool *pond.WorkerPool, viduc *video.UseCase, fuc *file.UseCase, vuc *view.UseCase, subuc *subscription.UseCase) *UseCase {
	siteUrl, ok := os.LookupEnv("SITE_URL")
	if !ok {
		log.Println("movie: SITE_URL not found in environment")
	}

	return &UseCase{
		MovieInteractor: &Repository{
			DB: db,
		},
		SiteUrl:             siteUrl,
		Pool:                pool,
		VideoUseCase:        viduc,
		FileUseCase:         fuc,
		ViewUseCase:         vuc,
		SubscriptionUseCase: subuc,
		MoviePool:           make(map[string]PoolItem),
		Mutex:               &sync.RWMutex{},
	}
}

func (uc *UseCase) Add(movieAddRequest *AddRequest) (*AddResponse, error) {
	// Limit max simultaneous uploads per user
	count, err := uc.MovieInteractor.GetWhereCount(map[string]interface{}{
		"user_id": movieAddRequest.UserId,
		"status":  StatusUploading,
	})
	if err != nil {
		return nil, err
	}

	if count >= 3 {
		return nil, errors.New("movie: too many movies")
	}

	movie := NewAddRequest(movieAddRequest)

	hasher := sha256.New()
	randomNumber := rand.Intn(1000)
	hasher.Write([]byte(time.Now().String() + strconv.Itoa(randomNumber)))
	hash := hasher.Sum(nil)

	encoded := base64.URLEncoding.EncodeToString(hash)
	encoded = strings.TrimRight(encoded, "=")

	movie.Code = encoded[:11]

	err = uc.MovieInteractor.Add(movie)
	if err != nil {
		return nil, err
	}

	return NewAddResponse(movie), nil
}

func (uc *UseCase) Delete(code string) error {
	movie, err := uc.MovieInteractor.GetSelectWhere(
		"id",
		map[string]interface{}{
			"code": code,
		},
	)
	if err != nil {
		return err
	}

	uc.Mutex.RLock()
	if poolItem, ok := uc.MoviePool[code]; ok {
		poolItem.Cancel()
	}
	uc.Mutex.RUnlock()

	err = uc.MovieInteractor.Delete(movie.ID)
	if err != nil {
		return err
	}

	go uc.FileUseCase.DeleteAllInPath("movies/" + code)

	return nil
}

func (uc *UseCase) DeleteMultiple(userId uint, movies *[]DeleteRequest) error {
	for _, movie := range *movies {
		if ok := uc.IsMovieOwner(userId, movie.Code); !ok {
			return errors.New("permission denied")
		}

		if err := uc.Delete(movie.Code); err != nil {
			return err
		}
	}

	return nil
}

func (uc *UseCase) UploadVideo(header *VideoUploadHeader, conn *websocket.Conn) error {
	movie, err := uc.MovieInteractor.Get(header.MovieCode)
	if err != nil {
		return errors.New("movie not found")
	}

	movieUpdateRequest := &VideoUpdateRequest{
		Code: header.MovieCode,
		Name: header.Filename,
	}
	rowsAffected, err := uc.UpdateVideo(movieUpdateRequest)
	if err != nil || rowsAffected == 0 {
		uc.Delete(movie.Code)
		return errors.New("failed to update video")
	}

	quality := video.GetQuality(1)

	tmpFilePath := filepath.Join("upload/movies", movie.Code, "resize")
	tmpFileName := quality.Code + ".mp4"
	tmpFile, err := uc.FileUseCase.WriteFileFromSocket(
		tmpFilePath,
		tmpFileName,
		[]string{"video/mp4", "video/avi", "video/webm"},
		header.Size,
		conn,
	)
	if err != nil {
		uc.Delete(movie.Code)
		return err
	}

	ctx, cancel := context.WithCancel(context.TODO())
	uc.Mutex.Lock()
	uc.MoviePool[movie.Code] = PoolItem{ctx, cancel}
	uc.Mutex.Unlock()

	savedVideo, err := uc.VideoUseCase.Save(ctx, tmpFile.Name(), tmpFileName, tmpFilePath, quality.ID)
	if err != nil {
		uc.Delete(movie.Code)
		return err
	}

	err = uc.MovieInteractor.AppendAssociation(&Movie{ID: movie.ID}, "Videos", savedVideo)
	if err != nil {
		uc.VideoUseCase.Delete(savedVideo)
		uc.Delete(movie.Code)
		return errors.New("failed to update video")
	}

	movie.Videos = append(movie.Videos, *savedVideo)

	uc.Pool.Submit(func() {
		uc.PostProcessVideo(ctx, *movie, tmpFile)
		uc.Mutex.Lock()
		delete(uc.MoviePool, movie.Code)
		uc.Mutex.Unlock()
	})

	return nil
}

func (uc *UseCase) RetryVideoPostProcess() {
	movies, err := uc.MovieInteractor.GetWhereMultiple(
		map[string]interface{}{"status": StatusUploading},
		&pagination.Pagination{
			Limit:  -1,
			Offset: -1,
		},
		"",
	)
	if err != nil {
		return
	}

	for _, movie := range *movies {
		var tmpVideo *video.Video
		for _, video := range movie.Videos {
			if video.Quality.ID == 1 {
				tmpVideo = &video
				break
			}
		}

		if tmpVideo == nil {
			uc.Delete(movie.Code)
			continue
		}

		tmpFile, err := os.Open(tmpVideo.File.Path)
		if err != nil {
			uc.Delete(movie.Code)
			continue
		}
		tmpFile.Close()

		movieCopy := movie

		ctx, cancel := context.WithCancel(context.TODO())
		uc.Mutex.RLock()
		uc.MoviePool[movie.Code] = PoolItem{ctx, cancel}
		uc.Mutex.RUnlock()

		uc.Pool.Submit(func() {
			uc.PostProcessVideo(ctx, movieCopy, tmpFile)
			uc.Mutex.Lock()
			delete(uc.MoviePool, movie.Code)
			uc.Mutex.Unlock()
		})
	}
}

func (uc *UseCase) PostProcessVideo(ctx context.Context, movie Movie, tmpFile *os.File) error {
	resizedVideoPath := filepath.Join("upload/movies", movie.Code, "resize")
	thumbsPath := filepath.Join("upload/movies", movie.Code, "thumbs")

	// Thumbs
	if err := uc.CreateThumbnails(ctx, movie, thumbsPath, tmpFile); err != nil {
		uc.Delete(movie.Code)
		return err
	}

	// Resize
	if err := uc.CreateResizedVideos(ctx, movie, resizedVideoPath, tmpFile); err != nil {
		uc.Delete(movie.Code)
		return err
	}

	os.Remove(tmpFile.Name())
	for _, video := range movie.Videos {
		if video.Quality.ID == 1 {
			uc.VideoUseCase.Delete(&video)
			break
		}
	}
	/*movie, err := uc.MovieInteractor.Get(movie.Code)
	if err == nil {

	}*/

	if uc.FileUseCase.GetSaveType() == file.SaveTypeInternal {
		os.RemoveAll(filepath.Join("upload/movies", movie.Code))
	}

	return nil
}

func (uc *UseCase) CreateThumbnails(ctx context.Context, movie Movie, thumbsPath string, tmpFile *os.File) error {
	if movie.WebVtt != nil {
		return nil
	}

	frameDuration := 10
	err := ffmpegthumbs.SplitVideoToThumbnails(tmpFile.Name(), thumbsPath, frameDuration)
	if err != nil {
		return errors.New("movie thumbnails: failed to create thumbnails")
	}

	thumbsWebvttPath := "/api/file/"
	imagesFilePath := make([]string, 0)
	var preview *file.File
	var previewWebp *file.File

	items, _ := os.ReadDir(thumbsPath)
	defaultPreviewPos := 1
	if len(items) > 2 {
		defaultPreviewPos = len(items) / 2
	}
	for i, item := range items {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if item.IsDir() {
				continue
			}

			imageFilePath := filepath.Join(thumbsPath, item.Name())
			savedImageFile, _ := uc.FileUseCase.CreateFromPath(imageFilePath, item.Name(), thumbsPath, "public")
			imagesFilePath = append(imagesFilePath, filepath.Join(thumbsWebvttPath, savedImageFile.Name))

			if defaultPreviewPos == i+1 {
				preview = savedImageFile
			}
		}
	}

	if len(imagesFilePath) == 0 {
		return errors.New("movie thumbnails: no thumbnails")
	}

	if preview != nil {
		previewWebp, _ = uc.FileUseCase.ImageToWebp(
			preview.FullPath, preview.Name, thumbsPath,
		)
	}

	var savedVttFile *file.File
	videoDuration, _ := ffmpegthumbs.GetVideoDuration(tmpFile.Name())
	vttFile, err := webvtt.CreateFromFilePaths(imagesFilePath, thumbsPath, videoDuration, frameDuration)
	if err != nil {
		return err
	}
	savedVttFile, _ = uc.FileUseCase.CreateFromPath(vttFile.Name(), "thumbs.vtt", thumbsPath, "public")

	movieUpdateRequest := &VideoUpdateRequest{
		Code:               movie.Code,
		DefaultPreview:     preview,
		DefaultPreviewWebp: previewWebp,
		WebVtt:             savedVttFile,
	}
	rowsAffected, err := uc.UpdateVideo(movieUpdateRequest)
	if err != nil || rowsAffected == 0 {
		return errors.New("failed to update video")
	}

	return nil
}

func (uc *UseCase) CreateResizedVideos(ctx context.Context, movie Movie, resizedVideoPath string, tmpFile *os.File) error {
	var resizedWebmPath string
	var savedVideo *video.Video
	var qualitiesIds []uint
	for _, movieVideo := range movie.Videos {
		qualitiesIds = append(qualitiesIds, movieVideo.Quality.ID)
	}

	_, videoHeight, _ := ffmpegthumbs.GetVideoSize(tmpFile.Name())

	for _, quality := range video.SupportedQualities {
		if slices.Contains(qualitiesIds, quality.ID) || videoHeight <= quality.Settings.MinHeight {
			continue
		}
		err := quality.Process(ctx, tmpFile.Name(), resizedVideoPath)
		if err != nil {
			return err
		}

		resizedWebmPath = filepath.Join(resizedVideoPath, quality.Code+".mp4")
		savedVideo, err = uc.VideoUseCase.Save(ctx, resizedWebmPath, quality.Code+".mp4", resizedVideoPath, quality.ID)
		if err != nil {
			return err
		}

		err = uc.MovieInteractor.AppendAssociation(&Movie{ID: movie.ID}, "Videos", savedVideo)
		if err != nil {
			return errors.New("failed to update video")
		}
	}

	movieUpdateRequest := &VideoUpdateRequest{
		Status: StatusReady,
		Code:   movie.Code,
	}
	rowsAffected, err := uc.UpdateVideo(movieUpdateRequest)
	if err != nil || rowsAffected == 0 {
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
	var selectQuery []string

	if utf8.RuneCountInString(movie.Name) > 130 {
		return errors.New("movie name too long")
	}
	if utf8.RuneCountInString(movie.Description) > 5000 {
		return errors.New("movie description too long")
	}

	if utf8.RuneCountInString(movie.Name) > 0 {
		selectQuery = append(selectQuery, "Name")
	}
	if utf8.RuneCountInString(movie.Description) > 0 {
		selectQuery = append(selectQuery, "Description")
	}

	if movie.Category.ID > 0 {
		selectQuery = append(selectQuery, "Category")
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
		isCorrectType, previewFileType := uc.FileUseCase.VerifyFileType(buff, []string{
			"image/jpeg", "image/png", "image/webp", "image/gif",
		})
		if !isCorrectType {
			return errors.New("invalid file type")
		}

		previewSavePath := filepath.Join("movies", movie.Code, "thumbs")
		preview, err := uc.FileUseCase.Create(
			movie.Preview, movie.PreviewHeader.Filename, previewSavePath, "public",
		)
		if err != nil {
			return err
		}

		if previewFileType == "image/gif" {
			movieRequest.PreviewWebpId = &preview.ID
		} else {
			previewWebp, err := uc.FileUseCase.ImageToWebp(preview.FullPath, preview.Name, "upload/"+previewSavePath)
			if err == nil {
				movieRequest.PreviewWebpId = &previewWebp.ID
			}
		}

		movieRequest.PreviewId = &preview.ID
		selectQuery = append(selectQuery, "PreviewId", "PreviewWebpId")
		uc.RemovePreview(movie.Code)
	} else if movie.RemovePreview {
		uc.RemovePreview(movie.Code)
	}

	rowsAffected, err := uc.MovieInteractor.UpdatesSelectWhere(
		movieRequest,
		selectQuery,
		map[string]interface{}{"code": movie.Code, "user_id": userId},
	)
	if err != nil {
		return err
	} else if rowsAffected == 0 {
		return errors.New("movie not found")
	}

	return nil
}

func (uc *UseCase) RemovePreview(code string) error {
	movie, err := uc.MovieInteractor.GetPreloadWhere(
		[]string{"Preview", "PreviewWebp"},
		map[string]interface{}{"code": code},
	)
	if err != nil {
		return err
	}

	if movie.Preview != nil {
		uc.FileUseCase.Delete(movie.Preview.Name)
	}
	if movie.PreviewWebp != nil {
		uc.FileUseCase.Delete(movie.PreviewWebp.Name)
	}

	return nil
}

func (uc *UseCase) UpdatePublishStatus(userId uint, movie *UpdatePublishStatusRequest) (int64, error) {
	movieRequest := NewUpdatePublishStatusRequest(movie)
	return uc.MovieInteractor.UpdatesSelectWhere(
		movieRequest,
		[]string{"is_published"},
		map[string]interface{}{"code": movie.Code, "user_id": userId},
	)
}

func (uc *UseCase) Get(userId *uint, code string) (*GetResponse, error) {
	movie, err := uc.MovieInteractor.Get(code)
	if err != nil {
		return nil, err
	}

	if movie.IsPublished || (userId != nil && movie.UserId == *userId) {
		response := NewGetResponse(movie)

		return response, nil
	}

	return nil, errors.New("not allowed")
}

func (uc *UseCase) GetPublic(userId *uint, code string, userIp net.IP) (*GetResponse, error) {
	movie, err := uc.MovieInteractor.Get(code)
	if err != nil {
		return nil, err
	}

	if movie.IsPublished || (userId != nil && movie.UserId == *userId) {
		response := NewGetResponse(movie)

		viewsCount, err := uc.ViewUseCase.GetCount(movie.ID)
		if err == nil {
			response.Views = viewsCount
		}

		if userId != nil {
			subscription, _ := uc.SubscriptionUseCase.Get(*userId, movie.UserId)
			if subscription != nil {
				response.Subscribed = ptr.Bool(subscription.ID > 0)
			}
		}

		view, err := uc.ViewUseCase.Add(movie.ID, userId, userIp)
		if err == nil {
			uc.MovieInteractor.AppendAssociation(&Movie{ID: movie.ID}, "Views", view)
		}

		return response, nil
	}

	return nil, errors.New("not allowed")
}

func (uc *UseCase) IsMovieOwner(userId uint, code string) bool {
	_, err := uc.MovieInteractor.GetSelectWhere(
		"id",
		map[string]interface{}{
			"user_id": userId,
			"code":    code,
		},
	)
	if err != nil {
		return false
	}

	return true
}

func (uc *UseCase) CheckMovieAccess(userId *uint, code string) bool {
	movie, err := uc.MovieInteractor.GetSelectWhere(
		[]string{"IsPublished", "UserId"},
		map[string]interface{}{"code": code},
	)
	if err != nil {
		return false
	}

	if movie.IsPublished || (userId != nil && movie.UserId == *userId) {
		return true
	}

	return false
}

func (uc *UseCase) CheckByUser(userId uint, code string) bool {
	_, err := uc.MovieInteractor.GetSelectWhere(
		"id",
		map[string]interface{}{
			"user_id": userId,
			"code":    code,
			"status":  StatusUploading,
		},
	)
	if err != nil {
		return false
	}

	return true
}

func (uc *UseCase) GetMultipleByUserId(userId uint, pagination *pagination.Pagination, sorting *sorting.Sort) ([]*GetForUserResponse, error) {
	if pagination.Limit > 20 || pagination.Limit == -1 {
		pagination.Limit = 20
	}

	order := ""
	if slices.Contains([]string{"created_at"}, sorting.SortBy) {
		order = fmt.Sprintf("%s %s", sorting.SortBy, sorting.SortVal)
	} else {
		order = "created_at desc"
	}

	movies, err := uc.MovieInteractor.GetMultipleByUserId(userId, pagination, order)
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
	movie, err := uc.MovieInteractor.GetWhere(map[string]interface{}{
		"user_id": userId,
		"code":    code,
	})
	if err != nil {
		return nil, err
	}

	return NewGetForUserResponse(movie), nil
}

func (uc *UseCase) GetMultiple(where interface{}, pagination *pagination.Pagination, sorting *sorting.Sort) ([]*GetResponse, error) {
	if pagination.Limit > 20 || pagination.Limit == -1 {
		pagination.Limit = 100
	}

	order := ""
	if slices.Contains([]string{"created_at"}, sorting.SortBy) {
		order = fmt.Sprintf("%s %s", sorting.SortBy, sorting.SortVal)
	} else {
		order = "created_at desc"
	}

	movies, err := uc.MovieInteractor.GetPreloadWhereMultiple(
		[]string{"Preview", "PreviewWebp", "DefaultPreview", "DefaultPreviewWebp", "WebVtt", "User", "User.Picture"},
		where,
		pagination,
		order,
	)
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

	var moviesIds []uint
	for _, movie := range moviesPayload {
		moviesIds = append(moviesIds, movie.ID)
	}

	viewsCounts, err := uc.ViewUseCase.GetMultipleCount(moviesIds)
	if err == nil {
		for key, movie := range moviesPayload {
			if _, ok := viewsCounts[movie.ID]; ok {
				moviesPayload[key].Views = viewsCounts[movie.ID]
			}
		}
	}

	if sorting.SortBy == "views" {
		switch sorting.SortVal {
		case "desc":
			sort.Slice(moviesPayload, func(i, j int) bool {
				return moviesPayload[i].Views > moviesPayload[j].Views
			})
			break
		case "asc":
			sort.Slice(moviesPayload, func(i, j int) bool {
				return moviesPayload[i].Views < moviesPayload[j].Views
			})
			break
		}
	}

	return moviesPayload, nil
}

func (uc *UseCase) GetMultipleByChannel(channelId uint, pagination *pagination.Pagination, sorting *sorting.Sort) ([]*GetResponse, error) {
	return uc.GetMultiple(
		map[string]interface{}{"is_published": 1, "user_id": channelId},
		pagination,
		sorting,
	)
}

func (uc *UseCase) GetMultiplePublic(pagination *pagination.Pagination, sorting *sorting.Sort) ([]*GetResponse, error) {
	return uc.GetMultiple(map[string]interface{}{"is_published": 1}, pagination, sorting)
}

func (uc *UseCase) GetMultipleSubscribed(userId uint, pagination *pagination.Pagination) ([]*GetResponse, error) {
	if pagination.Limit > 20 || pagination.Limit == -1 {
		pagination.Limit = 20
	}

	subscriptions, err := uc.SubscriptionUseCase.GetAll(userId)
	if err != nil {
		return nil, err
	}

	var usersIds []uint
	for _, subscription := range subscriptions {
		usersIds = append(usersIds, subscription.ChannelID)
	}

	movies, err := uc.MovieInteractor.GetPreloadWhereMultiple(
		[]string{"Preview", "PreviewWebp", "DefaultPreview", "DefaultPreviewWebp", "WebVtt", "User", "User.Picture"},
		map[string]interface{}{"is_published": 1, "user_id": usersIds},
		pagination,
		"created_at desc",
	)
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

	var moviesIds []uint
	for _, movie := range moviesPayload {
		moviesIds = append(moviesIds, movie.ID)
	}

	viewsCounts, err := uc.ViewUseCase.GetMultipleCount(moviesIds)
	if err == nil {
		for key, movie := range moviesPayload {
			if _, ok := viewsCounts[movie.ID]; ok {
				moviesPayload[key].Views = viewsCounts[movie.ID]
			}
		}
	}

	return moviesPayload, nil
}

func (uc *UseCase) GetMovieDetailSeo(movieCode string, r *http.Request) (map[string]string, error) {
	movie, err := uc.Get(nil, movieCode)
	if err != nil {
		return nil, err
	}

	var moviePreview string
	if movie.Preview != nil {
		moviePreview = uc.SiteUrl + "/api/file/" + movie.Preview.Name
	} else if movie.DefaultPreview != nil {
		moviePreview = uc.SiteUrl + "/api/file/" + movie.DefaultPreview.Name
	}

	languageCode := language.GetLanguageCode(r)
	siteName, err := language.GetMessage("SITE_NAME", languageCode)
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"title":       movie.Name + " - " + siteName,
		"description": movie.Description,
		"image":       moviePreview,
	}, nil
}
