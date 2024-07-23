package ffmpegthumbs

import (
	"encoding/json"
	"fmt"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Probe struct {
	Streams []Stream `json:"streams"`
	Format  Format   `json:"format"`
}

type Format struct {
	Bitrate string `json:"bit_rate"`
}

type Stream struct {
	DurationTs int    `json:"duration_ts"`
	Duration   string `json:"duration"`
	RFrameRate string `json:"r_frame_rate"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
}

func SplitVideoToThumbnails(filePath, outputPath string, frameDuration int) error {
	err := os.MkdirAll(outputPath, os.ModePerm)
	if err != nil {
		return err
	}

	duration, err := GetVideoDuration(filePath)
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

func GetVideoDuration(filePath string) (int, error) {
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

func GetVideoSize(filePath string) (int, int, error) {
	probe := &Probe{}
	fileInfoJson, err := ffmpeg.Probe(filePath)
	if err != nil {
		return 0, 0, err
	}

	err = json.Unmarshal([]byte(fileInfoJson), &probe)
	if err != nil {
		return 0, 0, err
	}

	return probe.Streams[0].Width, probe.Streams[0].Height, nil
}

func GetVideoBitrate(filePath string) (int, error) {
	probe := &Probe{}
	fileInfoJson, err := ffmpeg.Probe(filePath)
	if err != nil {
		return 0, err
	}

	err = json.Unmarshal([]byte(fileInfoJson), &probe)
	if err != nil {
		return 0, err
	}

	bitrate, err := strconv.Atoi(probe.Format.Bitrate)
	if err != nil {
		return 0, err
	}

	return bitrate, nil
}

func Resize(height int, crf, speed, bitrate, filePath, outputPath, fileName string) error {
	err := os.MkdirAll(outputPath, os.ModePerm)
	if err != nil {
		return err
	}

	err = ffmpeg.
		Input(filePath).
		Filter("scale", ffmpeg.Args{fmt.Sprintf("-1:%d", height)}).
		Output(filepath.Join(outputPath, fileName+".webm"), ffmpeg.KwArgs{
			"map":   "0:a:0",
			"c:v":   "libvpx-vp9",
			"crf":   crf,
			"speed": speed,
			"b:v":   bitrate,
			"c:a":   "libopus",
		}).
		Silent(true).
		Run()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func ToWebm(filePath, crf, speed, bitrate, outputPath, fileName string) error {
	err := os.MkdirAll(outputPath, os.ModePerm)
	if err != nil {
		return err
	}

	err = ffmpeg.
		Input(filePath).
		Output(filepath.Join(outputPath, fileName+".webm"), ffmpeg.KwArgs{
			"c:v":   "libvpx-vp9",
			"crf":   crf,
			"speed": speed,
			"b:v":   bitrate,
			"c:a":   "libopus",
		}).
		Silent(true).
		Run()
	if err != nil {
		return err
	}

	return nil
}
