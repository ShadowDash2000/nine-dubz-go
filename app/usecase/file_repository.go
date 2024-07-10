package usecase

import (
	"github.com/gorilla/websocket"
	"nine-dubz/app/model"
)

type FileRepository interface {
	Add(file *model.File) (uint, error)
	Remove(id uint) error
	Save(file *model.File) error
	Updates(file *model.File) error
	Get(id uint) (*model.File, error)
	CopyTmpFile(uploadPath string, tmpFilePath string, header *model.UploadHeader) error
	WriteFileFromSocket(tmpPath string, fileTypes []string, header *model.UploadHeader, conn *websocket.Conn) (string, error)
}
