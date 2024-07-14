package file

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
	"net/http"
)

type Handler struct {
	FileUseCase *UseCase
}

func NewHandler(uc *UseCase) *Handler {
	return &Handler{
		FileUseCase: uc,
	}
}

func (h *Handler) StreamFile(w http.ResponseWriter, r *http.Request) {
	fileName := chi.URLParam(r, "fileName")
	requestRange := r.Header.Get("Range")

	buff, contentRange, contentLength, err := h.FileUseCase.StreamFile(fileName, requestRange)
	if err != nil {
		http.Error(w, "No such file", http.StatusBadRequest)
		return
	}

	w.Header().Set("Accept-Ranges", "bytes")
	if len(requestRange) > 0 {
		w.Header().Set("Content-Range", contentRange)
		w.Header().Set("Content-Length", contentLength)
	}
	w.Header().Set("Content-Type", "video/mp4")
	w.WriteHeader(http.StatusPartialContent)
	w.Write(buff)
}

func (h *Handler) SocketUpload(w http.ResponseWriter, r *http.Request, fileTypes []string) (*File, error) {
	conn, err := h.FileUseCase.UpgradeConnection(w, r)
	if err != nil {
		fmt.Println("Error upgrading to websocket")
		return nil, err
	}
	defer conn.Close()

	header := &UploadHeader{}
	messageType, message, err := conn.ReadMessage()
	if err != nil {
		fmt.Println("Error receiving websocket message:", err)
		return nil, err
	}

	if messageType != websocket.TextMessage {
		err = errors.New("invalid message received, expecting file name and length")
		fmt.Println(err)
		return nil, err
	}
	if err := json.Unmarshal(message, header); err != nil {
		fmt.Println("Error receiving file name and length: " + err.Error())
		return nil, err
	}
	if len(header.Filename) == 0 {
		err = errors.New("filename cannot be empty")
		return nil, err
	}
	if header.Size == 0 {
		err = errors.New("upload file is empty")
		return nil, err
	}

	file, err := h.FileUseCase.WriteFileFromSocket(fileTypes, header, conn)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (h *Handler) SocketVideoUpload(w http.ResponseWriter, r *http.Request) (*File, error) {
	return h.SocketUpload(w, r, []string{"video/mp4"})
}
