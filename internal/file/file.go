package file

import (
	"context"
	"io"
	"net/http"
	"nine-dubz/pkg/ffmpegthumbs"
	"nine-dubz/pkg/s3storage"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type UseCase struct {
	FileInteractor Interactor
	IsDev          bool
}

func New(db *gorm.DB) *UseCase {
	isDevStr, ok := os.LookupEnv("IS_DEV")
	if !ok {
		isDevStr = "false"
	}
	isDev, err := strconv.ParseBool(isDevStr)
	if err != nil {
		isDev = false
	}
	saveType, ok := os.LookupEnv("FILE_SAVE_TYPE")
	if !ok {
		saveType = "local"
	}

	return &UseCase{
		FileInteractor: &Repository{
			DB:        db,
			S3Storage: s3storage.NewS3Storage(),
			SaveType:  saveType,
		},
		IsDev: isDev,
	}
}

func (uc *UseCase) UpgradeConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upgrader := websocket.Upgrader{}

	if uc.IsDev {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	}

	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	return connection, nil
}

func (uc *UseCase) Create(file io.ReadSeeker, name, path string, fileType string) (*File, error) {
	return uc.FileInteractor.Create(file, name, path, fileType)
}

func (uc *UseCase) CreateMultipart(ctx context.Context, filePath, name, path, fileType string) (*File, error) {
	return uc.FileInteractor.CreateMultipart(ctx, filePath, name, path, fileType)
}

func (uc *UseCase) CreateFromPath(filePath, name, path, fileType string) (*File, error) {
	return uc.FileInteractor.CreateFromPath(filePath, name, path, fileType)
}

func (uc *UseCase) Get(name string) ([]byte, error) {
	return uc.FileInteractor.Get(name)
}

func (uc *UseCase) Stream(file *File, requestRange string) ([]byte, string, int, error) {
	return uc.FileInteractor.Stream(file, requestRange)
}

func (uc *UseCase) Delete(name string) error {
	return uc.FileInteractor.Delete(name)
}

func (uc *UseCase) DeleteMultiple(names []string) error {
	return uc.FileInteractor.DeleteMultiple(names)
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

func (uc *UseCase) ImageToWebp(imagePath, name, savePath string) (*File, error) {
	err := ffmpegthumbs.ToWebp(
		imagePath,
		savePath,
		name,
	)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(filepath.Join(savePath, name+".webp"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return uc.Create(file, name+".webp", savePath, "public")
}
