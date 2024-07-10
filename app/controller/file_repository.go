package controller

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"io"
	"net/http"
	"nine-dubz/app/model"
	"os"
	"path/filepath"
	"time"
)

type FileRepository struct {
	DB *gorm.DB
}

func (fr *FileRepository) Add(file *model.File) (uint, error) {
	result := fr.DB.Create(&file)

	return file.ID, result.Error
}

func (fr *FileRepository) Remove(id uint) error {
	result := fr.DB.Delete(&model.File{}, id)

	return result.Error
}

func (fr *FileRepository) Save(file *model.File) error {
	result := fr.DB.Save(&file)

	return result.Error
}

func (fr *FileRepository) Updates(file *model.File) error {
	result := fr.DB.Updates(&file)

	return result.Error
}

func (fr *FileRepository) Get(id uint) (*model.File, error) {
	file := &model.File{}
	result := fr.DB.First(&file, id)

	return file, result.Error
}

func VerifyFileType(buff []byte, types []string) (bool, string) {
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

func (fr *FileRepository) CopyTmpFile(uploadPath string, tmpFilePath string, header *model.UploadHeader) error {
	tmpFile, err := os.Open(tmpFilePath)
	if err != nil {
		fmt.Println(err)
		return err
	}

	err = os.MkdirAll(uploadPath, os.ModePerm)
	if err != nil {
		return err
	}

	newFileName := time.Now().UnixNano()
	filePath := fmt.Sprintf(uploadPath+"/%d%s", newFileName, filepath.Ext(header.Filename))
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = fr.Add(&model.File{
		Name:         newFileName,
		OriginalName: header.Filename,
		Path:         filePath,
	})
	if err != nil {
		return err
	}

	_, err = io.Copy(f, tmpFile)
	if err != nil {
		return err
	}

	return nil
}

const (
	UploadStatusNext     int = 0
	UploadStatusError    int = 1
	UploadStatusComplete int = 2
)

func (fr *FileRepository) WriteFileFromSocket(tmpPath string, fileTypes []string, header *model.UploadHeader, conn *websocket.Conn) (string, error) {
	tmpFile, err := os.CreateTemp(tmpPath, "websocket_upload_")
	if err != nil {
		fmt.Println("Could not create temp file: " + err.Error())
		return "", err
	}
	defer tmpFile.Close()

	bytesRead := 0
	isCorrectType := false

	conn.WriteJSON(&model.UploadStatus{
		Status: UploadStatusNext,
	})

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			conn.WriteJSON(&model.UploadStatus{
				Status: UploadStatusError,
				Error:  "failed to receive message: " + err.Error(),
			})
			return "", err
		}

		if !isCorrectType {
			isCorrectType, _ = VerifyFileType(message, fileTypes)
			if !isCorrectType {
				conn.WriteJSON(&model.UploadStatus{
					Status: UploadStatusError,
					Error:  "file type not supported",
				})
				return "", errors.New("file type not supported")
			}
		}

		if bytesRead > 2<<30 {
			conn.WriteJSON(&model.UploadStatus{
				Status: UploadStatusError,
				Error:  err.Error(),
			})
			fmt.Println("File is too large")
			return "", errors.New("file is too large")
		}

		if messageType == websocket.CloseMessage {
			conn.WriteJSON(&model.UploadStatus{
				Status: UploadStatusComplete,
				Error:  "upload canceled",
			})
			return "", errors.New("upload canceled")
		}

		if messageType != websocket.BinaryMessage {
			conn.WriteJSON(&model.UploadStatus{
				Status: UploadStatusError,
				Error:  "invalid file block received",
			})
			return "", errors.New("invalid file block received")
		}

		tmpFile.Write(message)

		bytesRead += len(message)
		if bytesRead == header.Size {
			tmpFile.Close()
			break
		}

		conn.WriteJSON(&model.UploadStatus{
			Status: UploadStatusNext,
		})
	}

	conn.WriteJSON(&model.UploadStatus{
		Status: UploadStatusComplete,
	})

	return tmpFile.Name(), nil
}
