package controller

import "nine-dubz/app/usecase"

type VideoController struct {
	VideoInteractor usecase.VideoInteractor
}

func NewVideoController() *VideoController {
	return &VideoController{
		VideoInteractor: usecase.VideoInteractor{
			VideoRepository: &VideoRepository{},
		},
	}
}

func (vc *VideoController) SplitVideoToThumbnails(filePath string, outputPath string) error {
	return vc.VideoInteractor.SplitVideoToThumbnails(filePath, outputPath)
}

func (vc *VideoController) Resize(filePath string, outputPath string, fileName string) error {
	return vc.VideoInteractor.Resize(filePath, outputPath, fileName)
}
