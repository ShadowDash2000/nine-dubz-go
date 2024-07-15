package movie

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/gorilla/websocket"
	"net/http"
	"nine-dubz/internal/file"
	"nine-dubz/internal/pagination"
	"nine-dubz/internal/user"
)

type Handler struct {
	MovieUseCase *UseCase
	UserHandler  *user.Handler
	FileUseCase  *file.UseCase
}

func NewHandler(uc *UseCase, uh *user.Handler, fuc *file.UseCase) *Handler {
	return &Handler{
		MovieUseCase: uc,
		UserHandler:  uh,
		FileUseCase:  fuc,
	}
}

func (h *Handler) AddHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId")
	if userId == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	movieAddRequest := &AddRequest{
		UserId: userId.(uint),
	}
	movieAddResponse, err := h.MovieUseCase.Add(movieAddRequest)
	if err != nil {
		http.Error(w, "Can't add movie", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, movieAddResponse)
}

func (h *Handler) UploadVideoHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId")
	if userId == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := h.FileUseCase.UpgradeConnection(w, r)
	if err != nil {
		http.Error(w, "Can't upgrade connection", http.StatusInternalServerError)
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

	if ok := h.MovieUseCase.CheckByUser(userId.(uint), header.MovieCode); !ok {
		conn.WriteJSON(&file.UploadStatus{
			Status: file.UploadStatusError,
			Error:  "Permission denied",
		})
		return
	}

	file, err := h.FileUseCase.WriteFileFromSocket([]string{"video/mp4"}, header.Size, header.Filename, conn)
	if err != nil {
		h.MovieUseCase.Delete(userId.(uint), header.MovieCode)
		return
	}

	movieUpdateRequest := &UpdateRequest{
		Code:  header.MovieCode,
		Video: file,
	}
	h.MovieUseCase.Update(movieUpdateRequest)
}

func (h *Handler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	movieUpdateRequest := &UpdateRequest{}
	if err := json.NewDecoder(r.Body).Decode(movieUpdateRequest); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err := h.MovieUseCase.Update(movieUpdateRequest)
	if err != nil {
		http.Error(w, "Can't update movie", http.StatusBadRequest)
		return
	}

	render.JSON(w, r, struct {
		IsSuccess bool `json:"isSuccess"`
	}{true})
}

func (h *Handler) GetHandler(w http.ResponseWriter, r *http.Request) {
	movieCode := chi.URLParam(r, "movieCode")

	movie, err := h.MovieUseCase.Get(movieCode)
	if err != nil {
		http.Error(w, "Movie not found", http.StatusNotFound)
		return
	}

	render.JSON(w, r, movie)
}

func (h *Handler) GetMultipleForUserHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId")
	if userId == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	pagination := r.Context().Value("pagination").(*pagination.Pagination)

	moviesResponse, err := h.MovieUseCase.GetMultipleByUserId(userId.(uint), pagination)
	if err != nil {
		http.Error(w, "Movies not found", http.StatusNotFound)
		return
	}

	render.JSON(w, r, moviesResponse)
}

func (h *Handler) GetMultipleHandler(w http.ResponseWriter, r *http.Request) {
	pagination := r.Context().Value("pagination").(*pagination.Pagination)

	moviesResponse, err := h.MovieUseCase.GetMultiple(pagination)
	if err != nil {
		http.Error(w, "Movies not found", http.StatusNotFound)
		return
	}

	render.JSON(w, r, moviesResponse)
}

func (h *Handler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId")
	if userId == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	movieCode := chi.URLParam(r, "movieCode")

	if err := h.MovieUseCase.Delete(userId.(uint), movieCode); err != nil {
		http.Error(w, "Movie not found", http.StatusNotFound)
		return
	}
}

func (h *Handler) DeleteMultipleHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId")
	if userId == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	moviesDeleteRequest := &[]DeleteRequest{}
	if err := json.NewDecoder(r.Body).Decode(moviesDeleteRequest); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if err := h.MovieUseCase.DeleteMultiple(userId.(uint), moviesDeleteRequest); err != nil {
		http.Error(w, "Movie not found", http.StatusNotFound)
	}
}

func (h *Handler) UpdatePublishStatusHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("userId")
	if userId == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	movieUpdatePublishStatusRequest := &UpdatePublishStatusRequest{}
	if err := json.NewDecoder(r.Body).Decode(movieUpdatePublishStatusRequest); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if err := h.MovieUseCase.UpdatePublishStatus(userId.(uint), movieUpdatePublishStatusRequest); err != nil {
		http.Error(w, "Movie not found", http.StatusNotFound)
		return
	}
}
