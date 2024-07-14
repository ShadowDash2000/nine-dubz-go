package ffmpegthumbs

import (
	"encoding/json"
	"fmt"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"os"
	"strconv"
	"strings"
	"time"
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

func (vr *FfmpegThumbs) SplitVideoToThumbnails(filePath string, outputPath string) error {
	err := os.MkdirAll(outputPath, os.ModePerm)
	if err != nil {
		return err
	}

	probe := &Probe{}
	fileInfoJson, err := ffmpeg.Probe(filePath)
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(fileInfoJson), &probe)
	if err != nil {
		return err
	}

	duration, err := strconv.Atoi(strings.Split(probe.Streams[0].Duration, ".")[0])
	if err != nil {
		return err
	}

	file, err := os.Create(outputPath + "/thumbs.vtt")
	if err != nil {
		return err
	}
	defer file.Close()

	file.Write([]byte("WEBVTT\n\n"))

	i := 0
	for second := 1; second < duration; second = second + 10 {
		frameStart := time.Time{}
		frameStart = frameStart.Add(time.Duration(second) * time.Second)
		frameEnd := time.Time{}
		frameEnd = frameEnd.Add(time.Duration(second+10) * time.Second)
		timeFormat := "15:04:05.000"

		timeString := frameStart.Format(timeFormat) + " --> " + frameEnd.Format(timeFormat)
		imagePath := fmt.Sprintf("%s/img%d.jpg", outputPath, i)

		file.Write([]byte(timeString + "\n"))
		file.Write([]byte(imagePath + "\n\n"))

		ffmpeg.
			Input(filePath, ffmpeg.KwArgs{"ss": second}).
			Output(imagePath, ffmpeg.KwArgs{"vframes": "1", "q:v": "10"}).
			Silent(true).
			Run()

		i = i + 1
	}

	return nil
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
