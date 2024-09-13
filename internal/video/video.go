package video

import (
	"gorm.io/gorm"
	"nine-dubz/internal/file"
	"nine-dubz/pkg/ffmpegthumbs"
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

func (uc *UseCase) Save(filePath string, qualityId uint) (*Video, error) {
	savedFile, err := uc.FileUseCase.CreateFromPath(filePath, "private")
	if err != nil {
		return nil, err
	}
	width, height, _ := ffmpegthumbs.GetVideoSize(filePath)

	video := &Video{
		Width:   width,
		Height:  height,
		File:    savedFile,
		Quality: Quality{ID: qualityId},
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
