package file

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/gorilla/websocket"
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"io"
	"math/rand"
	"net/http"
	"nine-dubz/pkg/s3storage"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Repository struct {
	DB        *gorm.DB
	S3Storage *s3storage.S3Storage
}

func (fr *Repository) Create(file io.ReadSeeker, name, path string, size int64, fileType string) (*File, error) {
	timeNow := time.Now().UnixNano()
	randomNumber := rand.Intn(1000)

	newFileName := fmt.Sprintf("%d%d", timeNow, randomNumber)
	extension := filepath.Ext(name)

	_, err := fr.S3Storage.PutObject(file, newFileName+extension, path)
	if err != nil {
		return nil, err
	}

	savedFile := &File{
		Name:         newFileName,
		Extension:    extension,
		OriginalName: name,
		Size:         size,
		Path:         path,
		Type:         fileType,
	}
	result := fr.DB.Create(&savedFile)
	if result.Error != nil {
		fr.S3Storage.DeleteObject(newFileName, path)
		return nil, err
	}

	return savedFile, result.Error
}

func (fr *Repository) CreateMultipart(ctx context.Context, file io.ReadSeeker, name, path string, size int64, fileType string) (*File, error) {
	timeNow := time.Now().UnixNano()
	randomNumber := rand.Intn(1000)

	newFileName := fmt.Sprintf("%d%d", timeNow, randomNumber)
	extension := filepath.Ext(name)

	_, err := fr.S3Storage.MultipartUpload(ctx, file, size, newFileName+extension, path)
	if err != nil {
		return nil, err
	}

	savedFile := &File{
		Name:         newFileName,
		Extension:    extension,
		OriginalName: name,
		Size:         size,
		Path:         path,
		Type:         fileType,
	}
	result := fr.DB.Create(&savedFile)
	if result.Error != nil {
		fr.S3Storage.DeleteObject(newFileName, path)
		return nil, err
	}

	return savedFile, result.Error
}

func (fr *Repository) Get(name string) ([]byte, error) {
	file, err := fr.GetWhere(map[string]interface{}{
		"name": name,
		"type": "public",
	})
	if err != nil {
		return nil, err
	}

	output, err := fr.S3Storage.GetObject(name+file.Extension, file.Path)
	if err != nil {
		return nil, err
	}

	buff, err := io.ReadAll(output.Body)
	if err != nil {
		return nil, err
	}

	return buff, nil
}

func (fr *Repository) Stream(name, path, requestRange string) ([]byte, string, int64, error) {
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
	output, err := fr.S3Storage.GetRangeObject(name, path, "bytes="+contentRange)
	if err != nil {
		return nil, "", 0, err
	}

	contentRange = aws.ToString(output.ContentRange)
	contentLength := aws.ToInt64(output.ContentLength)

	buff, _ = io.ReadAll(output.Body)

	return buff, contentRange, contentLength, nil
}

func (fr *Repository) Delete(name string) error {
	file, err := fr.GetWhere(map[string]interface{}{"name": name})
	if err != nil {
		return err
	}

	result := fr.DB.Unscoped().Delete(&File{}, "name = ?", name)

	if result.Error == nil {
		fr.S3Storage.DeleteObject(name, file.Path)
	}

	return result.Error
}

func (fr *Repository) DeleteMultiple(names []string, path string) error {
	result := fr.DB.Unscoped().Delete(&File{}, "name IN ?", names)

	if result.Error == nil {
		fr.S3Storage.DeleteObjects(names, path)
	}

	return result.Error
}

func (fr *Repository) DeleteAllInPath(path string) error {
	result := fr.DB.Unscoped().Delete(&File{}, "path = ?", path)

	if result.Error == nil {
		fr.S3Storage.DeleteAllInPrefix(path)
	}

	return result.Error
}

func (fr *Repository) Updates(file *File) error {
	result := fr.DB.Updates(&file)

	return result.Error
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

func (fr *Repository) WriteFileFromSocket(tmpPath, fileName string, fileTypes []string, fileSize int, maxChunkSize int, conn *websocket.Conn) (*os.File, error) {
	err := os.MkdirAll(tmpPath, os.ModePerm)
	if err != nil {
		return nil, err
	}

	tmpFile, err := os.Create(filepath.Join(tmpPath, fileName))
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

		if len(message) > maxChunkSize {
			tmpFile.Close()
			os.Remove(tmpFile.Name())

			conn.WriteJSON(&UploadStatus{
				Status: UploadStatusError,
				Error:  "chunk too large",
			})
			return nil, errors.New("chunk too large")
		}

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
