package file

import (
	"github.com/gorilla/websocket"
	"mime/multipart"
)

type Interactor interface {
	Add(file *File) (*File, error)
	Remove(id uint) error
	Save(file *File) error
	Updates(file *File) error
	Get(id uint) (*File, error)
	VerifyFileType(buff []byte, types []string) (bool, string)
	CopyTmpFile(uploadPath string, tmpFilePath string, header *UploadHeader) (*File, error)
	WriteFileFromSocket(tmpPath string, fileTypes []string, header *UploadHeader, conn *websocket.Conn) (string, error)
	SaveFile(path string, fileName string, file multipart.File) (*File, error)
}
