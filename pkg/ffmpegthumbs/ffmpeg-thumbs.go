package ffmpegthumbs

import (
	"encoding/json"
	"fmt"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"golang.org/x/net/context"
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
	CodecType  string `json:"codec_type"`
	DurationTs int    `json:"duration_ts"`
	Duration   string `json:"duration"`
	RFrameRate string `json:"r_frame_rate"`
	Width      int    `json:"width"`
	Height     int    `json:"height"`
	Bitrate    string `json:"bit_rate"`
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

	var duration string
	for _, stream := range probe.Streams {
		if stream.CodecType == "video" {
			duration = strings.Split(stream.Duration, ".")[0]
			break
		}
	}

	if duration == "" {
		return 0, fmt.Errorf("ffmpeg: duration is empty")
	}

	return strconv.Atoi(duration)
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

	var width, height int
	for _, stream := range probe.Streams {
		if stream.CodecType == "video" {
			width = stream.Width
			height = stream.Height
			break
		}
	}

	if width == 0 || height == 0 {
		return 0, 0, fmt.Errorf("ffmpeg: width or height is empty")
	}

	return width, height, nil
}

func GetVideoBitrate(filePath string) (string, error) {
	return GetBitrate(filePath, "video")
}

func GetAudioBitrate(filePath string) (string, error) {
	return GetBitrate(filePath, "audio")
}

func GetBitrate(filePath, codecType string) (string, error) {
	probe := &Probe{}
	fileInfoJson, err := ffmpeg.Probe(filePath)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal([]byte(fileInfoJson), &probe)
	if err != nil {
		return "", err
	}

	var bitrate string
	for _, stream := range probe.Streams {
		if stream.CodecType == codecType {
			bitrate = stream.Bitrate
			break
		}
	}

	if bitrate == "" {
		return "", fmt.Errorf("ffmpeg: bitrate is empty")
	}

	return bitrate, nil
}

func Resize(ctx context.Context, height int, crf, speed, videoBitrate, audioBitrate, filePath, outputPath, fileName string) error {
	err := os.MkdirAll(outputPath, os.ModePerm)
	if err != nil {
		return err
	}

	stream := ffmpeg.
		Input(filePath).
		Filter("scale", ffmpeg.Args{fmt.Sprintf("-2:%d", height)}).
		Output(filepath.Join(outputPath, fileName+".mp4"), ffmpeg.KwArgs{
			"map":   "0:a:0",
			"c:v":   "libx264",
			"crf":   crf,
			"speed": speed,
			"b:v":   videoBitrate,
			"c:a":   "libopus",
			"b:a":   audioBitrate,
		}).
		Silent(true).
		OverWriteOutput()

	stream.Context, _ = context.WithCancel(ctx)
	err = stream.Run()
	if err != nil {
		return err
	}

	return nil
}

func ToWebm(ctx context.Context, filePath, crf, speed, bitrate, outputPath, fileName string) error {
	err := os.MkdirAll(outputPath, os.ModePerm)
	if err != nil {
		return err
	}

	stream := ffmpeg.
		Input(filePath).
		Output(filepath.Join(outputPath, fileName+".mp4"), ffmpeg.KwArgs{
			"c:v":   "libx264",
			"crf":   crf,
			"speed": speed,
			"b:v":   bitrate,
			"c:a":   "libopus",
		}).
		Silent(true).
		OverWriteOutput()

	stream.Context, _ = context.WithCancel(ctx)
	err = stream.Run()
	if err != nil {
		return err
	}

	return nil
}

func ToWebp(filePath, outputPath, fileName string) error {
	err := os.MkdirAll(outputPath, os.ModePerm)
	if err != nil {
		return err
	}

	err = ffmpeg.
		Input(filePath).
		Output(filepath.Join(outputPath, fileName+".webp"), ffmpeg.KwArgs{
			"c:v": "libwebp",
		}).
		Silent(true).
		OverWriteOutput().
		Run()
	if err != nil {
		return err
	}

	return nil
}
