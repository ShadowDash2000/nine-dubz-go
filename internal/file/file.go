package file

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"io"
	"net/http"
	"nine-dubz/pkg/s3storage"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type UseCase struct {
	FileInteractor Interactor
	S3Storage      *s3storage.S3Storage
}

func New(db *gorm.DB) *UseCase {
	return &UseCase{
		FileInteractor: &Repository{
			DB: db,
		},
		S3Storage: s3storage.NewS3Storage(),
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

func (uc *UseCase) StreamFile(fileName, requestRange string) ([]byte, string, string, error) {
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
	output, err := uc.S3Storage.GetRangeObject(fileName, "bytes="+contentRange)
	if err != nil {
		return nil, "", "", err
	}

	contentRange = aws.ToString(output.ContentRange)
	contentLength := strconv.Itoa(int(aws.ToInt64(output.ContentLength)))

	buff, _ = io.ReadAll(output.Body)

	return buff, contentRange, contentLength, nil
}

func (uc *UseCase) GetFile(fileName string) ([]byte, error) {
	_, err := uc.FileInteractor.GetWhere(map[string]interface{}{
		"name": fileName,
		"type": "public",
	})
	if err != nil {
		return nil, err
	}

	output, err := uc.S3Storage.GetObject(fileName)
	if err != nil {
		return nil, err
	}

	buff, err := io.ReadAll(output.Body)
	if err != nil {
		return nil, err
	}

	return buff, nil
}

func (uc *UseCase) SaveFile(file io.Reader, fileName string, fileSize int64, fileType string) (*File, error) {
	timeNow := time.Now().UnixNano()
	newFileName := strconv.Itoa(int(timeNow))
	extension := filepath.Ext(fileName)

	_, err := uc.S3Storage.PutObject(file, newFileName)
	if err != nil {
		return nil, err
	}

	savedFile, err := uc.FileInteractor.Add(&File{
		Name:         newFileName,
		Extension:    extension,
		OriginalName: fileName,
		Size:         fileSize,
		Type:         fileType,
	})
	if err != nil {
		uc.S3Storage.DeleteObject(newFileName)
		return nil, err
	}

	return savedFile, nil
}

func (uc *UseCase) VerifyFileType(buff []byte, types []string) (bool, string) {
	return uc.FileInteractor.VerifyFileType(buff, types)
}

func (uc *UseCase) WriteFileFromSocket(fileTypes []string, fileSize int, fileName string, conn *websocket.Conn) (*File, *os.File, error) {
	tmpFile, err := uc.FileInteractor.WriteFileFromSocket("upload/tmp", fileTypes, fileSize, 1024*1024, conn)
	if err != nil {
		return nil, nil, err
	}

	conn.Close()

	timeNow := time.Now().UnixNano()
	newFileName := strconv.Itoa(int(timeNow))
	extension := filepath.Ext(fileName)

	tmpFile, _ = os.Open(tmpFile.Name())
	defer tmpFile.Close()

	_, err = uc.S3Storage.PutObject(tmpFile, newFileName)
	if err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return nil, nil, err
	}

	savedFile, err := uc.FileInteractor.Add(&File{
		Name:         newFileName,
		Extension:    extension,
		OriginalName: fileName,
		Size:         int64(fileSize),
		Type:         "private",
	})
	if err != nil {
		os.Remove(tmpFile.Name())
		uc.S3Storage.DeleteObject(newFileName)
		return nil, nil, err
	}

	return savedFile, tmpFile, nil
}
