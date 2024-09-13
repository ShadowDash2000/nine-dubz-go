package file

import (
	"github.com/gorilla/websocket"
	"io"
	"os"
)

type Interactor interface {
	Create(file io.ReadSeeker, name, path string, fileType string) (*File, error)
	CreateFromPath(path string, fileType string) (*File, error)
	Get(name string) ([]byte, error)
	Stream(file *File, requestRange string) ([]byte, string, int, error)
	Delete(name string) error
	DeleteMultiple(names []string) error
	DeleteAllInPath(path string) error
	Updates(file *File) error
	GetWhere(where map[string]interface{}) (*File, error)
	GetWhereMultiple(where map[string]interface{}) ([]File, error)
	VerifyFileType(buff []byte, types []string) (bool, string)
	WriteFileFromSocket(tmpPath, fileName string, fileTypes []string, fileSize int, maxChunkSize int, conn *websocket.Conn) (*os.File, error)
}
