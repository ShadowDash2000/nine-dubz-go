package file

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"io"
	"io/fs"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Repository struct {
	DB *gorm.DB
}

const SaveFolderPrefix = "upload/"

func (fr *Repository) Create(file io.ReadSeeker, name, path string, fileType string) (*File, error) {
	timeNow := time.Now().UnixNano()
	randomNumber := rand.Intn(1000)

	newFileName := fmt.Sprintf("%d%d", timeNow, randomNumber)
	extension := filepath.Ext(name)

	err := os.MkdirAll(filepath.Join(SaveFolderPrefix, path), os.ModePerm)
	if err != nil {
		return nil, err
	}

	f, err := os.Create(filepath.Join(SaveFolderPrefix, path, newFileName+extension))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fileInfo, err := os.Stat(f.Name())
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(f, file)
	if err != nil {
		return nil, err
	}

	savedFile := &File{
		Name:         newFileName,
		Extension:    extension,
		OriginalName: name,
		Size:         fileInfo.Size(),
		Path:         f.Name(),
		Type:         fileType,
	}
	result := fr.DB.Create(&savedFile)
	if result.Error != nil {
		return nil, result.Error
	}

	return savedFile, nil
}

func (fr *Repository) CreateFromPath(path string, fileType string) (*File, error) {
	timeNow := time.Now().UnixNano()
	randomNumber := rand.Intn(1000)
	newFileName := fmt.Sprintf("%d%d", timeNow, randomNumber)

	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	extension := filepath.Ext(path)

	savedFile := &File{
		Name:         newFileName,
		Extension:    extension,
		OriginalName: fileInfo.Name(),
		Size:         fileInfo.Size(),
		Path:         path,
		Type:         fileType,
	}
	result := fr.DB.Create(&savedFile)
	if result.Error != nil {
		return nil, result.Error
	}

	return savedFile, nil
}

func (fr *Repository) Get(name string) ([]byte, error) {
	file, err := fr.GetWhere(map[string]interface{}{
		"name": name,
		"type": "public",
	})
	if err != nil {
		return nil, err
	}

	f, err := os.Open(file.Path)

	buff, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return buff, nil
}

func (fr *Repository) Stream(file *File, requestRange string) ([]byte, string, int, error) {
	f, err := os.Open(file.Path)
	if err != nil {
		return nil, "", 0, err
	}
	defer f.Close()

	var offset int64
	if len(requestRange) > 0 {
		requestRange = strings.Replace(requestRange, "bytes=", "", -1)
		requestRange = strings.Replace(requestRange, "-", "", -1)
		offset, _ = strconv.ParseInt(requestRange, 10, 64)
	} else {
		offset = 0
	}

	_, err = f.Seek(offset, 0)
	if err != nil {
		return nil, "", 0, err
	}

	buffSize := int64(1024 * 1024 * 5)
	if buffSize+offset > file.Size {
		buffSize = file.Size - offset
	}
	buff := make([]byte, buffSize)
	contentLength, err := f.Read(buff)
	if err != nil {
		return nil, "", 0, err
	}

	contentRange := "bytes " + strconv.FormatInt(offset, 10) + "-" +
		strconv.Itoa(int(offset)+contentLength-1) + "/" +
		strconv.FormatInt(file.Size, 10)

	return buff, contentRange, contentLength, nil
}

func (fr *Repository) Delete(name string) error {
	file, err := fr.GetWhere(map[string]interface{}{"name": name})
	if err != nil {
		return err
	}

	result := fr.DB.Unscoped().Delete(&File{}, "name = ?", name)
	if result.Error != nil {
		return result.Error
	}

	err = os.Remove(file.Path)
	if err != nil {
		return err
	}

	return result.Error
}

func (fr *Repository) DeleteMultiple(names []string) error {
	files, err := fr.GetWhereMultiple(map[string]interface{}{"name": names})
	if err != nil {
		return err
	}

	result := fr.DB.Unscoped().Delete(&File{}, "name IN ?", names)

	if result.Error != nil {
		return result.Error
	}

	for _, file := range files {
		err = os.Remove(file.Path)
		if err != nil {
			return err
		}
	}

	return result.Error
}

func (fr *Repository) DeleteAllInPath(path string) error {
	var paths []string
	err := filepath.Walk(filepath.Join(SaveFolderPrefix, path), func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	})

	result := fr.DB.Unscoped().Delete(&File{}, map[string]interface{}{"path": paths})
	if result.Error != nil {
		return result.Error
	}

	err = os.RemoveAll(filepath.Join(SaveFolderPrefix, path))
	if err != nil {
		return err
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

func (fr *Repository) GetWhereMultiple(where map[string]interface{}) ([]File, error) {
	var file []File
	result := fr.DB.Where(where).Find(&file)

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
