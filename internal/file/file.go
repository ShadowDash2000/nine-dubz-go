package file

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"io"
	"net/http"
	"nine-dubz/pkg/s3storage"
	"os"
	"path/filepath"
)

type UseCase struct {
	FileInteractor Interactor
}

func New(db *gorm.DB) *UseCase {
	return &UseCase{
		FileInteractor: &Repository{
			DB:        db,
			S3Storage: s3storage.NewS3Storage(),
		},
	}
}

func (uc *UseCase) UpgradeConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upgrader := websocket.Upgrader{}
	/**
	TODO added for testing
	REMOVE IN FUTURE
	*/
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	return connection, nil
}

func (uc *UseCase) Create(file io.ReadSeeker, name, path string, size int64, fileType string) (*File, error) {
	return uc.FileInteractor.Create(file, name, path, size, fileType)
}

func (uc *UseCase) CreateMultipart(ctx context.Context, file io.ReadSeeker, name, path string, size int64, fileType string) (*File, error) {
	return uc.FileInteractor.CreateMultipart(ctx, file, name, path, size, fileType)
}

func (uc *UseCase) Get(name string) ([]byte, error) {
	return uc.FileInteractor.Get(name)
}

func (uc *UseCase) Stream(name, path, requestRange string) ([]byte, string, int64, error) {
	return uc.FileInteractor.Stream(name, path, requestRange)
}

func (uc *UseCase) Delete(name string) error {
	return uc.FileInteractor.Delete(name)
}

func (uc *UseCase) DeleteMultiple(names []string, path string) error {
	return uc.FileInteractor.DeleteMultiple(names, path)
}

func (uc *UseCase) DeleteAllInPath(path string) error {
	return uc.FileInteractor.DeleteAllInPath(path)
}

func (uc *UseCase) VerifyFileType(buff []byte, types []string) (bool, string) {
	return uc.FileInteractor.VerifyFileType(buff, types)
}

func (uc *UseCase) WriteFileFromSocket(filePath, fileName string, fileTypes []string, fileSize int, conn *websocket.Conn) (*os.File, error) {
	tmpFile, err := uc.FileInteractor.WriteFileFromSocket(filePath, fileName, fileTypes, fileSize, 1024*1024, conn)
	if err != nil {
		return nil, err
	}

	conn.Close()

	tmpFile, _ = os.Open(tmpFile.Name())
	defer tmpFile.Close()

	return tmpFile, nil
}

func (uc *UseCase) DownloadFile(pathTo, name, pathFrom string, file *File) (*os.File, error) {
	err := os.MkdirAll(pathTo, os.ModePerm)
	if err != nil {
		return nil, err
	}
	tmpFile, err := os.Create(filepath.Join(pathTo, name))
	if err != nil {
		return nil, err
	}
	defer tmpFile.Close()

	var currentByte int64
	for {
		if currentByte >= file.Size {
			break
		}

		requestRange := fmt.Sprintf("bytes=%d-", currentByte)
		buff, _, contentLength, err := uc.FileInteractor.Stream(file.Name+file.Extension, pathFrom, requestRange)
		if err != nil {
			tmpFile.Close()
			os.Remove(tmpFile.Name())
			return nil, err
		}

		tmpFile.Write(buff)

		currentByte = currentByte + contentLength
	}

	if currentByte < file.Size {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return nil, errors.New("file: file was corrupted while downloading")
	}

	return tmpFile, nil
}
