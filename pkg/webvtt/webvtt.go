package webvtt

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

type WebVTT struct{}

func (wvtt *WebVTT) CreateFromFolder(folderPath, outputPath string, videoDuration, frameDuration int) error {
	file, err := os.Create(filepath.Join(outputPath, "thumbs.vtt"))
	if err != nil {
		return err
	}
	defer file.Close()

	file.Write([]byte("WEBVTT\n\n"))

	items, _ := os.ReadDir(folderPath)
	second := 0
	for _, item := range items {
		if item.IsDir() {
			continue
		}
		frameStart := time.Time{}
		frameStart = frameStart.Add(time.Duration(second) * time.Second)
		frameEnd := time.Time{}
		frameEnd = frameEnd.Add(time.Duration(second+frameDuration) * time.Second)
		timeFormat := "15:04:05.000"

		timeString := frameStart.Format(timeFormat) + " --> " + frameEnd.Format(timeFormat)

		file.Write([]byte(timeString + "\n"))
		file.Write([]byte(filepath.Join(folderPath, item.Name()) + "\n\n"))

		second = second + frameDuration

		if second >= videoDuration {
			break
		}
	}

	return nil
}

func (wvtt *WebVTT) CreateFromFilePaths(filePaths []string, outputPath string, videoDuration, frameDuration int) (*os.File, error) {
	file, err := os.Create(filepath.Join(outputPath, "thumbs.vtt"))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	file.Write([]byte("WEBVTT\n\n"))

	second := 0
	for _, filePath := range filePaths {
		frameStart := time.Time{}
		frameStart = frameStart.Add(time.Duration(second) * time.Second)
		frameEnd := time.Time{}
		frameEnd = frameEnd.Add(time.Duration(second+frameDuration) * time.Second)
		timeFormat := "15:04:05.000"

		timeString := frameStart.Format(timeFormat) + " --> " + frameEnd.Format(timeFormat)

		file.Write([]byte(timeString + "\n"))
		file.Write([]byte(strings.ReplaceAll(filePath, "\\", "/") + "\n\n"))

		second = second + frameDuration

		if second >= videoDuration {
			break
		}
	}

	return file, nil
}
