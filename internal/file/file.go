package file

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"io"
	"mime/multipart"
	"net/http"
	"nine-dubz/pkg/ffmpegthumbs"
	"nine-dubz/pkg/s3storage"
	"os"
	"strconv"
	"strings"
)

type UseCase struct {
	FileInteractor Interactor
	S3Storage      *s3storage.S3Storage
	FfmpegThumbs   *ffmpegthumbs.FfmpegThumbs
}

func New(db *gorm.DB) *UseCase {
	return &UseCase{
		FileInteractor: &Repository{
			DB: db,
		},
		S3Storage:    s3storage.NewS3Storage(),
		FfmpegThumbs: &ffmpegthumbs.FfmpegThumbs{},
	}
}

func (uc *UseCase) UpgradeConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upgrader := websocket.Upgrader{}
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
	output, err := uc.S3Storage.GetObject(fileName, "bytes="+contentRange)
	if err != nil {
		return nil, "", "", err
	}

	contentRange = aws.ToString(output.ContentRange)
	contentLength := strconv.Itoa(int(aws.ToInt64(output.ContentLength)))

	buff, _ = io.ReadAll(output.Body)

	return buff, contentRange, contentLength, nil
}

func (uc *UseCase) SaveFile(path string, fileName string, file multipart.File) (*File, error) {
	return uc.FileInteractor.SaveFile(path, fileName, file)
}

func (uc *UseCase) VerifyFileType(buff []byte, types []string) (bool, string) {
	return uc.FileInteractor.VerifyFileType(buff, types)
}

func (uc *UseCase) WriteFileFromSocket(fileTypes []string, header *UploadHeader, conn *websocket.Conn) (*File, error) {
	tmpFile, err := uc.FileInteractor.WriteFileFromSocket("upload/tmp", fileTypes, header, conn)
	if err != nil {
		return nil, err
	}

	conn.Close()

	file, err := uc.FileInteractor.CopyTmpFile("upload/video", tmpFile, header)
	if err != nil {
		return nil, err
	}

	fileReader, err := os.Open(file.Path)
	if err != nil {
		return nil, err
	}
	defer fileReader.Close()

	_, err = uc.S3Storage.PutObject(fileReader, file.Name)
	if err != nil {
		uc.FileInteractor.Remove(file.ID)
		os.Remove(file.Path)
		return nil, err
	}

	go func() {
		err = uc.FfmpegThumbs.SplitVideoToThumbnails(file.Path, "upload/thumbs/"+file.Name)
		if err != nil {
			fmt.Println(err)
		}

		os.Remove(file.Path)
	}()

	return file, nil
}
