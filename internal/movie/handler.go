package movie

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/gorilla/websocket"
	"net/http"
	"nine-dubz/internal/file"
	"nine-dubz/internal/pagination"
	"nine-dubz/internal/response"
	"nine-dubz/internal/token"
	"nine-dubz/internal/user"
	"nine-dubz/pkg/tokenauthorize"
	"strconv"
)

type Handler struct {
	MovieUseCase   *UseCase
	UserHandler    *user.Handler
	FileUseCase    *file.UseCase
	TokenAuthorize *tokenauthorize.TokenAuthorize
	TokenUseCase   *token.UseCase
}

func NewHandler(uc *UseCase, uh *user.Handler, fuc *file.UseCase, ta *tokenauthorize.TokenAuthorize, tuc *token.UseCase) *Handler {
	return &Handler{
		MovieUseCase:   uc,
		UserHandler:    uh,
		FileUseCase:    fuc,
		TokenAuthorize: ta,
		TokenUseCase:   tuc,
	}
}

func (h *Handler) AddHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId")
	if userId == "" {
		response.RenderError(w, r, http.StatusUnauthorized, "Unauthorized")
		return
	}

	movieAddRequest := &AddRequest{
		UserId: userId.(uint),
	}
	movieAddResponse, err := h.MovieUseCase.Add(movieAddRequest)
	if err != nil {
		response.RenderError(w, r, http.StatusInternalServerError, "Can't add movie")
		return
	}

	render.JSON(w, r, movieAddResponse)
}

func (h *Handler) UploadVideoHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := h.FileUseCase.UpgradeConnection(w, r)
	if err != nil {
		response.RenderError(w, r, http.StatusInternalServerError, "Can't upgrade connection")
		return
	}
	defer conn.Close()

	messageType, message, err := conn.ReadMessage()
	if err != nil {
		conn.WriteJSON(&file.UploadStatus{
			Status: file.UploadStatusError,
			Error:  "Error receiving websocket message",
		})
		return
	}

	if messageType != websocket.TextMessage {
		conn.WriteJSON(&file.UploadStatus{
			Status: file.UploadStatusError,
			Error:  "Invalid message received, expecting file name, length and movie code",
		})
		return
	}

	header := &VideoUploadHeader{}
	if err = json.Unmarshal(message, header); err != nil {
		conn.WriteJSON(&file.UploadStatus{
			Status: file.UploadStatusError,
			Error:  "Error receiving upload header",
		})
		return
	}
	if len(header.Filename) == 0 {
		conn.WriteJSON(&file.UploadStatus{
			Status: file.UploadStatusError,
			Error:  "File name cannot be empty",
		})
		return
	}
	if header.Size == 0 {
		conn.WriteJSON(&file.UploadStatus{
			Status: file.UploadStatusError,
			Error:  "File size cannot be zero",
		})
		return
	}
	if header.MovieCode == "" {
		conn.WriteJSON(&file.UploadStatus{
			Status: file.UploadStatusError,
			Error:  "No movie code",
		})
		return
	}
	if header.Token == "" {
		conn.WriteJSON(&file.UploadStatus{
			Status: file.UploadStatusError,
			Error:  "No token",
		})
		return
	}

	if _, err = h.TokenAuthorize.VerifyToken(header.Token); err != nil {
		conn.WriteJSON(&file.UploadStatus{
			Status: file.UploadStatusError,
			Error:  "Invalid token",
		})
		return
	}

	userId, err := h.TokenUseCase.GetUserIdByToken(header.Token)
	if err != nil {
		conn.WriteJSON(&file.UploadStatus{
			Status: file.UploadStatusError,
			Error:  "User not found",
		})
		return
	}

	if ok := h.MovieUseCase.CheckByUser(userId, header.MovieCode); !ok {
		conn.WriteJSON(&file.UploadStatus{
			Status: file.UploadStatusError,
			Error:  "Permission denied",
		})
		return
	}

	err = h.MovieUseCase.SaveVideo(userId, header, conn)
	if err != nil {
		fmt.Println(err)
	}
}

func (h *Handler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId")
	if userId == "" {
		response.RenderError(w, r, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if err := r.ParseMultipartForm(6 << 20); err != nil {
		response.RenderError(w, r, http.StatusBadRequest, "Failed to parse form data")
		return
	}

	movieUpdateRequest := &UpdateRequest{}
	movieUpdateRequest.Code = chi.URLParam(r, "movieCode")
	movieUpdateRequest.Name = r.PostForm.Get("name")
	movieUpdateRequest.Description = r.PostForm.Get("description")

	isPublished, err := strconv.ParseBool(r.PostForm.Get("isPublished"))
	if err == nil {
		movieUpdateRequest.IsPublished = isPublished
	}

	file, fileHeader, err := r.FormFile("preview")
	if err == nil {
		movieUpdateRequest.Preview = file
		movieUpdateRequest.PreviewHeader = fileHeader
	}

	err = h.MovieUseCase.UpdateByUserId(userId.(uint), movieUpdateRequest)
	if err != nil {
		response.RenderError(w, r, http.StatusBadRequest, "Can't update movie: "+err.Error())
		return
	}

	render.JSON(w, r, struct {
		IsSuccess bool `json:"isSuccess"`
	}{true})
}

func (h *Handler) GetHandler(w http.ResponseWriter, r *http.Request) {
	movieCode := chi.URLParam(r, "movieCode")
	userId := r.Context().Value("userId")
	if userId == nil {
		userId = uint(0)
	}

	movie, err := h.MovieUseCase.Get(userId.(uint), movieCode)
	if err != nil {
		response.RenderError(w, r, http.StatusNotFound, "Movie not found")
		return
	}

	render.JSON(w, r, movie)
}

func (h *Handler) GetForUserHandler(w http.ResponseWriter, r *http.Request) {
	movieCode := chi.URLParam(r, "movieCode")
	userId := r.Context().Value("userId")
	if userId == nil {
		userId = uint(0)
	}

	movie, err := h.MovieUseCase.GetForUser(userId.(uint), movieCode)
	if err != nil {
		response.RenderError(w, r, http.StatusNotFound, "Movie not found")
		return
	}

	render.JSON(w, r, movie)
}

func (h *Handler) GetMultipleForUserHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId")
	if userId == "" {
		response.RenderError(w, r, http.StatusUnauthorized, "Unauthorized")
		return
	}
	pagination := r.Context().Value("pagination").(*pagination.Pagination)

	moviesResponse, err := h.MovieUseCase.GetMultipleByUserId(userId.(uint), pagination)
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, make([]struct{}, 0))
		return
	}

	if len(moviesResponse) > 0 {
		render.JSON(w, r, moviesResponse)
	} else {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, make([]struct{}, 0))
	}
}

func (h *Handler) GetMultipleHandler(w http.ResponseWriter, r *http.Request) {
	pagination := r.Context().Value("pagination").(*pagination.Pagination)

	moviesResponse, err := h.MovieUseCase.GetMultiple(pagination)
	if err != nil {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, make([]struct{}, 0))
		return
	}

	if len(moviesResponse) > 0 {
		render.JSON(w, r, moviesResponse)
	} else {
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, make([]struct{}, 0))
	}
}

func (h *Handler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId")
	if userId == "" {
		response.RenderError(w, r, http.StatusUnauthorized, "Unauthorized")
		return
	}
	movieCode := chi.URLParam(r, "movieCode")

	if err := h.MovieUseCase.Delete(userId.(uint), movieCode); err != nil {
		response.RenderError(w, r, http.StatusNotFound, "Movie not found")
		return
	}

	response.RenderSuccess(w, r, http.StatusOK, "")
}

func (h *Handler) DeleteMultipleHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId")
	if userId == "" {
		response.RenderError(w, r, http.StatusUnauthorized, "Unauthorized")
		return
	}

	moviesDeleteRequest := &[]DeleteRequest{}
	if err := json.NewDecoder(r.Body).Decode(moviesDeleteRequest); err != nil {
		response.RenderError(w, r, http.StatusBadRequest, "Can't parse fields")
		return
	}

	if err := h.MovieUseCase.DeleteMultiple(userId.(uint), moviesDeleteRequest); err != nil {
		response.RenderError(w, r, http.StatusNotFound, "Movie not found")
	}

	response.RenderSuccess(w, r, http.StatusOK, "")
}

func (h *Handler) UpdatePublishStatusHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId")
	if userId == "" {
		response.RenderError(w, r, http.StatusUnauthorized, "Unauthorized")
		return
	}

	movieUpdatePublishStatusRequest := &UpdatePublishStatusRequest{}
	if err := json.NewDecoder(r.Body).Decode(movieUpdatePublishStatusRequest); err != nil {
		response.RenderError(w, r, http.StatusBadRequest, "Can't parse fields")
		return
	}

	rowsAffected, err := h.MovieUseCase.UpdatePublishStatus(userId.(uint), movieUpdatePublishStatusRequest)
	if err != nil || rowsAffected == 0 {
		response.RenderError(w, r, http.StatusNotFound, "Movie not found")
		return
	}
}

func (h *Handler) StreamFile(w http.ResponseWriter, r *http.Request) {
	movieCode := chi.URLParam(r, "movieCode")
	quality := r.URL.Query().Get("q")
	requestRange := r.Header.Get("Range")
	userId := r.Context().Value("userId")
	if userId == nil {
		userId = uint(0)
	}

	movie, err := h.MovieUseCase.CheckMovieAccess(userId.(uint), movieCode)
	if err != nil {
		response.RenderError(w, r, http.StatusNotFound, "Movie not found")
		return
	}

	file := movie.Video
	switch quality {
	case "0":
		file = movie.VideoShakal
		break
	case "360":
		file = movie.Video360
		break
	case "480":
		file = movie.Video480
		break
	case "720":
		file = movie.Video720
		break
	}

	if file == nil {
		response.RenderError(w, r, http.StatusNotFound, "No such quality")
		return
	}

	buff, contentRange, contentLength, err := h.FileUseCase.StreamFile(file.Name, requestRange)
	if err != nil {
		response.RenderError(w, r, http.StatusNotFound, "File not found")
		return
	}

	w.Header().Set("Accept-Ranges", "bytes")
	if len(requestRange) > 0 {
		w.Header().Set("Content-Range", contentRange)
		w.Header().Set("Content-Length", strconv.FormatInt(contentLength, 10))
	}
	//w.Header().Set("Content-Type", "video/mp4")
	w.WriteHeader(http.StatusPartialContent)
	w.Write(buff)
}
