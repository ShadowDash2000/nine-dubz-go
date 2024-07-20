package file

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"net/http"
	"os"
)

type Repository struct {
	DB *gorm.DB
}

func (fr *Repository) Add(file *File) (*File, error) {
	result := fr.DB.Create(&file)

	return file, result.Error
}

func (fr *Repository) Remove(id uint) error {
	result := fr.DB.Delete(&File{}, id)

	return result.Error
}

func (fr *Repository) Save(file *File) error {
	result := fr.DB.Save(&file)

	return result.Error
}

func (fr *Repository) Updates(file *File) error {
	result := fr.DB.Updates(&file)

	return result.Error
}

func (fr *Repository) Get(id uint) (*File, error) {
	file := &File{}
	result := fr.DB.First(&file, id)

	return file, result.Error
}

func (fr *Repository) GetWhere(where map[string]interface{}) (*File, error) {
	file := &File{}
	result := fr.DB.Where(where).First(&file)

	return file, result.Error
}

func (fr *Repository) VerifyFileType(buff []byte, types []string) (bool, string) {
	filetype := http.DetectContentType(buff)
	isCorrectType := false
	for _, t := range types {
		if filetype == t {
			isCorrectType = true
			break
		}
	}

	return isCorrectType, filetype
}

const (
	UploadStatusNextChunk int = 0
	UploadStatusError     int = 1
	UploadStatusComplete  int = 2
)

func (fr *Repository) WriteFileFromSocket(tmpPath string, fileTypes []string, fileSize int, conn *websocket.Conn) (*os.File, error) {
	err := os.MkdirAll(tmpPath, os.ModePerm)
	if err != nil {
		return nil, err
	}

	tmpFile, err := os.CreateTemp(tmpPath, "websocket_upload_")
	if err != nil {
		fmt.Println("Could not create temp file: " + err.Error())
		return nil, err
	}
	defer tmpFile.Close()

	bytesRead := 0
	isCorrectType := false

	conn.WriteJSON(&UploadStatus{
		Status: UploadStatusNextChunk,
	})

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			conn.WriteJSON(&UploadStatus{
				Status: UploadStatusError,
				Error:  "failed to receive message: " + err.Error(),
			})
			return nil, err
		}

		if !isCorrectType {
			isCorrectType, _ = fr.VerifyFileType(message, fileTypes)
			if !isCorrectType {
				conn.WriteJSON(&UploadStatus{
					Status: UploadStatusError,
					Error:  "file type not supported",
				})
				return nil, errors.New("file type not supported")
			}
		}

		if bytesRead > 8<<30 {
			conn.WriteJSON(&UploadStatus{
				Status: UploadStatusError,
				Error:  err.Error(),
			})
			fmt.Println("File is too large")
			return nil, errors.New("file is too large")
		}

		if messageType != websocket.BinaryMessage {
			tmpFile.Close()
			os.Remove(tmpFile.Name())

			if messageType == websocket.CloseMessage {
				conn.WriteJSON(&UploadStatus{
					Status: UploadStatusComplete,
					Error:  "upload canceled",
				})
				return nil, errors.New("upload canceled")
			}

			conn.WriteJSON(&UploadStatus{
				Status: UploadStatusError,
				Error:  "invalid file block received",
			})
			return nil, errors.New("invalid file block received")
		}

		/*
			TODO add max chunk size
		*/

		tmpFile.Write(message)

		bytesRead += len(message)
		if bytesRead == fileSize {
			tmpFile.Close()
			break
		} else if bytesRead > fileSize {
			tmpFile.Close()
			os.Remove(tmpFile.Name())

			conn.WriteJSON(&UploadStatus{
				Status: UploadStatusError,
				Error:  "read more than allowed file size",
			})
			return nil, errors.New("read more than allowed file size")
		}

		conn.WriteJSON(&UploadStatus{
			Status: UploadStatusNextChunk,
		})
	}

	conn.WriteJSON(&UploadStatus{
		Status: UploadStatusComplete,
	})

	return tmpFile, nil
}
