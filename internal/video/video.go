package video

import (
	"context"
	"nine-dubz/internal/file"
	"nine-dubz/pkg/ffmpegthumbs"

	"gorm.io/gorm"
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

func (uc *UseCase) Save(ctx context.Context, filePath, name, path string, qualityId uint) (*Video, error) {
	savedFile, err := uc.FileUseCase.CreateMultipart(ctx, filePath, name, path, "private")
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
