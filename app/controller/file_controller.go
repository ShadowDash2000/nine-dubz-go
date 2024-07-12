package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"io"
	"mime/multipart"
	"net/http"
	"nine-dubz/app/model"
	"nine-dubz/app/usecase"
	"os"
	"strconv"
	"strings"
)

type FileController struct {
	FileInteractor     usecase.FileInteractor
	SocketInteractor   usecase.SocketInteractor
	S3Controller       *S3Controller
	VideoController    *VideoController
	LanguageController *LanguageController
}

func NewFileController(db *gorm.DB, lc *LanguageController, vc *VideoController, sc *S3Controller) *FileController {
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
		VideoController:    vc,
		S3Controller:       sc,
		LanguageController: lc,
	}
}

func (fc *FileController) StreamFile(w http.ResponseWriter, r *http.Request) {
	fileName := chi.URLParam(r, "fileName")

	requestRange := r.Header.Get("Range")
	var off int
	if len(requestRange) > 0 {
		requestRange = strings.Replace(requestRange, "bytes=", "", -1)
		requestRange = strings.Replace(requestRange, "-", "", -1)
		off, _ = strconv.Atoi(requestRange)
	} else {
		off = 0
	}

	buff := make([]byte, 1024*1024*5)

	contentRange := strconv.Itoa(off) + "-" + strconv.Itoa(len(buff)+off)
	output, err := fc.S3Controller.Read(fileName, "bytes="+contentRange)
	if err != nil {
		http.Error(w, "No such file", http.StatusBadRequest)
		return
	}

	contentRange = aws.ToString(output.ContentRange)
	contentLength := strconv.Itoa(int(aws.ToInt64(output.ContentLength)))

	buff, _ = io.ReadAll(output.Body)

	w.Header().Set("Accept-Ranges", "bytes")
	if len(requestRange) > 0 {
		w.Header().Set("Content-Range", contentRange)
		w.Header().Set("Content-Length", contentLength)
	}
	w.Header().Set("Content-Type", "video/mp4")
	w.WriteHeader(http.StatusPartialContent)
	w.Write(buff)
}

func (fc *FileController) SocketUpload(w http.ResponseWriter, r *http.Request, fileTypes []string) (*model.File, error) {
	conn, err := fc.SocketInteractor.UpgradeConnection(w, r)
	if err != nil {
		fmt.Println("Error upgrading to websocket")
		return nil, err
	}
	defer conn.Close()

	header := &model.UploadHeader{}
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

	tmpFile, err := fc.FileInteractor.WriteFileFromSocket("upload/tmp", fileTypes, header, conn)
	if err != nil {
		fmt.Println(err)
		text, _ := fc.LanguageController.GetStringByCode(r, "WEBSOCKET_CANT_UPLOAD_FILE")
		render.Render(w, r, ErrInvalidRequest(err, http.StatusInternalServerError, text))
		return nil, err
	}

	conn.Close()

	file, err := fc.FileInteractor.CopyTmpFile("upload/video", tmpFile, header)
	if err != nil {
		text, _ := fc.LanguageController.GetStringByCode(r, "WEBSOCKET_FAILED_TO_WRITE_FILE")
		render.Render(w, r, ErrInvalidRequest(err, http.StatusInternalServerError, text))
		return nil, err
	}

	_, err = fc.S3Controller.Upload(file)
	if err != nil {
		fc.FileInteractor.Remove(file.ID)
		return nil, err
	}

	go func() {
		err = fc.VideoController.SplitVideoToThumbnails(file.Path, "upload/thumbs/"+file.Name)
		if err != nil {
			fmt.Println(err)
		}

		os.Remove(file.Path)
	}()

	return file, nil
}

func (fc *FileController) SaveFile(path string, fileName string, file multipart.File) (*model.File, error) {
	return fc.FileInteractor.SaveFile(path, fileName, file)
}

func (fc *FileController) VerifyFileType(buff []byte, types []string) (bool, string) {
	return fc.FileInteractor.VerifyFileType(buff, types)
}

func (fc *FileController) SocketVideoUpload(w http.ResponseWriter, r *http.Request) (*model.File, error) {
	return fc.SocketUpload(w, r, []string{"video/mp4"})
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
