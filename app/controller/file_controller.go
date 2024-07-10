package controller

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/render"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"net/http"
	"nine-dubz/app/model"
	"nine-dubz/app/usecase"
)

type FileController struct {
	FileInteractor     usecase.FileInteractor
	SocketInteractor   usecase.SocketInteractor
	LanguageController *LanguageController
}

func NewFileController(db *gorm.DB, lc *LanguageController) *FileController {
	return &FileController{
		FileInteractor: usecase.FileInteractor{
			FileRepository: &FileRepository{
				DB: db,
			},
		},
		SocketInteractor: usecase.SocketInteractor{
			SocketRepository: &SocketRepository{
				Clients: make(map[*websocket.Conn]bool),
			},
		},
		LanguageController: lc,
	}
}

func (fc *FileController) UploadHandler(w http.ResponseWriter, r *http.Request, fileTypes []string) {
	conn, err := fc.SocketInteractor.UpgradeConnection(w, r)
	if err != nil {
		text, _ := fc.LanguageController.GetStringByCode(r, "WEBSOCKET_CANT_UPGRADE")
		render.Render(w, r, ErrInvalidRequest(err, http.StatusInternalServerError, text))
		return
	}
	defer conn.Close()

	header := &model.UploadHeader{}
	messageType, message, err := conn.ReadMessage()
	if err != nil {
		fmt.Println("Error receiving websocket message:", err)
		return
	}

	if messageType != websocket.TextMessage {
		fmt.Println("Invalid message received, expecting file name and length")
		return
	}
	if err := json.Unmarshal(message, header); err != nil {
		fmt.Println("Error receiving file name and length: " + err.Error())
		return
	}
	if len(header.Filename) == 0 {
		fmt.Println("Filename cannot be empty")
		return
	}
	if header.Size == 0 {
		fmt.Println("Upload file is empty")
		return
	}

	tmpFile, err := fc.FileInteractor.WriteFileFromSocket("upload", fileTypes, header, conn)
	if err != nil {
		fmt.Println(err)
		text, _ := fc.LanguageController.GetStringByCode(r, "WEBSOCKET_CANT_UPLOAD_FILE")
		render.Render(w, r, ErrInvalidRequest(err, http.StatusInternalServerError, text))
		return
	}

	err = fc.FileInteractor.CopyTmpFile("upload", tmpFile, header)
	if err != nil {
		text, _ := fc.LanguageController.GetStringByCode(r, "WEBSOCKET_FAILED_TO_WRITE_FILE")
		render.Render(w, r, ErrInvalidRequest(err, http.StatusInternalServerError, text))
		return
	}
}

func (fc *FileController) UploadVideoHandler(w http.ResponseWriter, r *http.Request) {
	fc.UploadHandler(w, r, []string{"video/mp4"})
}

func (fc *FileController) MaxBodySize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 2<<30)
		if err := r.ParseForm(); err != nil {
			text, _ := fc.LanguageController.GetStringByCode(r, "REQUEST_MAX_SIZE_LIMIT")
			render.Render(w, r, ErrInvalidRequest(err, http.StatusBadRequest, text))
			return
		}

		next.ServeHTTP(w, r)
	})
}
