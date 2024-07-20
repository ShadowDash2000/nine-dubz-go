package file

import (
	"github.com/gorilla/websocket"
	"os"
)

type Interactor interface {
	Add(file *File) (*File, error)
	Remove(id uint) error
	Save(file *File) error
	Updates(file *File) error
	Get(id uint) (*File, error)
	GetWhere(where map[string]interface{}) (*File, error)
	VerifyFileType(buff []byte, types []string) (bool, string)
	WriteFileFromSocket(tmpPath string, fileTypes []string, fileSize int, conn *websocket.Conn) (*os.File, error)
}
