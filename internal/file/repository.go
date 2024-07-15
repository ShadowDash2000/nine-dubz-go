package file

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
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

func (fr *Repository) CopyTmpFile(uploadPath string, tmpFilePath string, fileName string) (*File, error) {
	tmpFile, err := os.Open(tmpFilePath)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer tmpFile.Close()

	err = os.MkdirAll(uploadPath, os.ModePerm)
	if err != nil {
		return nil, err
	}

	timeNow := time.Now().UnixNano()
	newFileName := strconv.Itoa(int(timeNow))
	extension := filepath.Ext(fileName)
	filePath := uploadPath + "/" + newFileName + extension
	f, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	size, err := io.Copy(f, tmpFile)
	if err != nil {
		return nil, err
	}

	tmpFile.Close()
	os.Remove(tmpFilePath)

	savedFile, err := fr.Add(&File{
		Name:         newFileName,
		Extension:    extension,
		OriginalName: fileName,
		Path:         filePath,
		Size:         size,
	})
	if err != nil {
		return nil, err
	}

	return savedFile, nil
}

const (
	UploadStatusNextChunk int = 0
	UploadStatusError     int = 1
	UploadStatusComplete  int = 2
)

func (fr *Repository) WriteFileFromSocket(tmpPath string, fileTypes []string, fileSize int, conn *websocket.Conn) (string, error) {
	err := os.MkdirAll(tmpPath, os.ModePerm)
	if err != nil {
		return "", err
	}

	tmpFile, err := os.CreateTemp(tmpPath, "websocket_upload_")
	if err != nil {
		fmt.Println("Could not create temp file: " + err.Error())
		return "", err
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
			return "", err
		}

		if !isCorrectType {
			isCorrectType, _ = fr.VerifyFileType(message, fileTypes)
			if !isCorrectType {
				conn.WriteJSON(&UploadStatus{
					Status: UploadStatusError,
					Error:  "file type not supported",
				})
				return "", errors.New("file type not supported")
			}
		}

		if bytesRead > 8<<30 {
			conn.WriteJSON(&UploadStatus{
				Status: UploadStatusError,
				Error:  err.Error(),
			})
			fmt.Println("File is too large")
			return "", errors.New("file is too large")
		}

		if messageType != websocket.BinaryMessage {
			tmpFile.Close()
			os.Remove(tmpFile.Name())

			if messageType == websocket.CloseMessage {
				conn.WriteJSON(&UploadStatus{
					Status: UploadStatusComplete,
					Error:  "upload canceled",
				})
				return "", errors.New("upload canceled")
			}

			conn.WriteJSON(&UploadStatus{
				Status: UploadStatusError,
				Error:  "invalid file block received",
			})
			return "", errors.New("invalid file block received")
		}

		tmpFile.Write(message)

		bytesRead += len(message)
		if bytesRead == fileSize {
			tmpFile.Close()
			break
		}

		conn.WriteJSON(&UploadStatus{
			Status: UploadStatusNextChunk,
		})
	}

	conn.WriteJSON(&UploadStatus{
		Status: UploadStatusComplete,
	})

	return tmpFile.Name(), nil
}

func (fr *Repository) SaveFile(path string, fileName string, file multipart.File) (*File, error) {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return nil, err
	}

	timeNow := time.Now().UnixNano()
	newFileName := strconv.Itoa(int(timeNow))
	extension := filepath.Ext(fileName)
	filePath := path + newFileName + extension
	f, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	size, err := io.Copy(f, file)
	if err != nil {
		return nil, err
	}

	savedFile, err := fr.Add(&File{
		Name:         newFileName,
		Extension:    extension,
		OriginalName: fileName,
		Path:         filePath,
		Size:         size,
	})
	if err != nil {
		return nil, err
	}

	return savedFile, nil
}
