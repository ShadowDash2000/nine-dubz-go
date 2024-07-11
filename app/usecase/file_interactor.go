package usecase

import (
	"github.com/gorilla/websocket"
	"mime/multipart"
	"nine-dubz/app/model"
)

type FileInteractor struct {
	FileRepository FileRepository
}

func (fi *FileInteractor) Add(file *model.File) (*model.File, error) {
	return fi.FileRepository.Add(file)
}

func (fi *FileInteractor) Remove(id uint) error {
	return fi.FileRepository.Remove(id)
}

func (fi *FileInteractor) Save(file *model.File) error {
	return fi.FileRepository.Save(file)
}

func (fi *FileInteractor) Updates(file *model.File) error {
	return fi.FileRepository.Updates(file)
}

func (fi *FileInteractor) Get(id uint) (*model.File, error) {
	return fi.FileRepository.Get(id)
}

func (fi *FileInteractor) VerifyFileType(buff []byte, types []string) (bool, string) {
	return fi.FileRepository.VerifyFileType(buff, types)
}

func (fi *FileInteractor) CopyTmpFile(uploadPath string, tmpFilePath string, header *model.UploadHeader) (*model.File, error) {
	return fi.FileRepository.CopyTmpFile(uploadPath, tmpFilePath, header)
}

func (fi *FileInteractor) WriteFileFromSocket(tmpPath string, fileTypes []string, header *model.UploadHeader, conn *websocket.Conn) (string, error) {
	return fi.FileRepository.WriteFileFromSocket(tmpPath, fileTypes, header, conn)
}

func (fi *FileInteractor) SaveFile(path string, fileName string, file multipart.File) (*model.File, error) {
	return fi.FileRepository.SaveFile(path, fileName, file)
}
