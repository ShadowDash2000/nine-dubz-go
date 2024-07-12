package usecase

type VideoInteractor struct {
	VideoRepository VideoRepository
}

func (vi *VideoInteractor) SplitVideoToThumbnails(filePath string, outputPath string) error {
	return vi.VideoRepository.SplitVideoToThumbnails(filePath, outputPath)
}

func (vi *VideoInteractor) Resize(filePath string, outputPath string, fileName string) error {
	return vi.VideoRepository.Resize(filePath, outputPath, fileName)
}
