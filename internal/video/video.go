package video

import (
	"gorm.io/gorm"
	"nine-dubz/internal/file"
	"nine-dubz/pkg/ffmpegthumbs"
	"os"
)

type UseCase struct {
	VideoInteractor Interactor
	FileUseCase     *file.UseCase
}

func New(db *gorm.DB, fuc *file.UseCase) *UseCase {
	return &UseCase{
		VideoInteractor: &Repository{
			DB: db,
		},
		FileUseCase: fuc,
	}
}

func (uc *UseCase) Save(filePath, pathTo string, qualityId uint) (*Video, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	fileInfo, _ := os.Stat(file.Name())
	savedFile, err := uc.FileUseCase.Create(file, fileInfo.Name(), pathTo, fileInfo.Size(), "private")
	if err != nil {
		return nil, err
	}
	file.Close()
	width, height, _ := ffmpegthumbs.GetVideoSize(filePath)

	video := &Video{
		Width:     width,
		Height:    height,
		File:      savedFile,
		QualityID: qualityId,
	}
	err = uc.VideoInteractor.Create(video)
	if err != nil {
		uc.FileUseCase.Delete(savedFile.Name)
		return nil, err
	}

	return video, nil
}

func (uc *UseCase) Delete(video *Video) error {
	err := uc.VideoInteractor.Delete(video.ID)
	if err != nil {
		return err
	}

	return uc.FileUseCase.Delete(video.File.Name)
}
