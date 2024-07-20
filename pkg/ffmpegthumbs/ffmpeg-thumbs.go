package ffmpegthumbs

import (
	"encoding/json"
	"fmt"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"os"
	"strconv"
	"strings"
)

type FfmpegThumbs struct{}

type Probe struct {
	Streams []Stream `json:"streams"`
}

type Stream struct {
	DurationTs int    `json:"duration_ts"`
	Duration   string `json:"duration"`
	RFrameRate string `json:"r_frame_rate"`
}

func (vr *FfmpegThumbs) SplitVideoToThumbnails(filePath, outputPath string, frameDuration int) error {
	err := os.MkdirAll(outputPath, os.ModePerm)
	if err != nil {
		return err
	}

	duration, err := vr.GetVideoDuration(filePath)
	if err != nil {
		return err
	}

	i := 0
	for second := 1; second < duration; second = second + frameDuration {
		imagePath := fmt.Sprintf("%s/img%06d.jpg", outputPath, i)

		ffmpeg.
			Input(filePath, ffmpeg.KwArgs{"ss": second}).
			Output(imagePath, ffmpeg.KwArgs{"vframes": "1", "q:v": frameDuration}).
			Silent(true).
			Run()

		i = i + 1
	}

	return nil
}

func (vr *FfmpegThumbs) GetVideoDuration(filePath string) (int, error) {
	probe := &Probe{}
	fileInfoJson, err := ffmpeg.Probe(filePath)
	if err != nil {
		return 0, err
	}

	err = json.Unmarshal([]byte(fileInfoJson), &probe)
	if err != nil {
		return 0, err
	}

	duration, err := strconv.Atoi(strings.Split(probe.Streams[0].Duration, ".")[0])
	if err != nil {
		return 0, err
	}

	return duration, nil
}

func (vr *FfmpegThumbs) Resize(filePath string, outputPath string, fileName string) error {
	err := os.MkdirAll(outputPath, os.ModePerm)
	if err != nil {
		return err
	}

	return ffmpeg.
		Input(filePath).
		Filter("scale", ffmpeg.Args{"-2:240"}).
		Output(outputPath+"/"+fileName, ffmpeg.KwArgs{"an": ""}).
		Silent(true).
		Run()
}
