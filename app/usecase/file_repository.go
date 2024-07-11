package usecase

import (
	"github.com/gorilla/websocket"
	"mime/multipart"
	"nine-dubz/app/model"
)

type FileRepository interface {
	Add(file *model.File) (*model.File, error)
	Remove(id uint) error
	Save(file *model.File) error
	Updates(file *model.File) error
	Get(id uint) (*model.File, error)
	VerifyFileType(buff []byte, types []string) (bool, string)
	CopyTmpFile(uploadPath string, tmpFilePath string, header *model.UploadHeader) (*model.File, error)
	WriteFileFromSocket(tmpPath string, fileTypes []string, header *model.UploadHeader, conn *websocket.Conn) (string, error)
	SaveFile(path string, fileName string, file multipart.File) (*model.File, error)
}
