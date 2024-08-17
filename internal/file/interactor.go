package file

import (
	"github.com/gorilla/websocket"
	"golang.org/x/net/context"
	"io"
	"os"
)

type Interactor interface {
	Create(file io.ReadSeeker, name, path string, size int64, fileType string) (*File, error)
	CreateMultipart(ctx context.Context, file io.ReadSeeker, name, path string, size int64, fileType string) (*File, error)
	Get(name string) ([]byte, error)
	Stream(name, path, requestRange string) ([]byte, string, int64, error)
	Delete(name string) error
	DeleteMultiple(names []string, path string) error
	DeleteAllInPath(path string) error
	Updates(file *File) error
	GetWhere(where map[string]interface{}) (*File, error)
	VerifyFileType(buff []byte, types []string) (bool, string)
	WriteFileFromSocket(tmpPath, fileName string, fileTypes []string, fileSize int, maxChunkSize int, conn *websocket.Conn) (*os.File, error)
}
