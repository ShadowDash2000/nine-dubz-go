package usecase

type VideoRepository interface {
	SplitVideoToThumbnails(filePath string, outputPath string) error
	Resize(filePath string, outputPath string, fileName string) error
}
