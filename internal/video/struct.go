package video

import (
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"nine-dubz/internal/file"
	"nine-dubz/pkg/ffmpegthumbs"
)

type Video struct {
	gorm.Model
	ID        uint
	QualityID uint   `json:"-"`
	Title     string `gorm:"-"`
	Width     int
	Height    int
	FileID    uint
	File      *file.File
}

type GetResponse struct {
	Quality *Quality   `json:"quality"`
	Width   int        `json:"width"`
	Height  int        `json:"height"`
	File    *file.File `json:"file"`
}

func NewGetResponse(video Video) *GetResponse {
	return &GetResponse{
		Quality: GetQuality(video.QualityID),
		Width:   video.Width,
		Height:  video.Height,
		File:    video.File,
	}
}

func NewGetResponseMultiple(videos []Video) []*GetResponse {
	response := make([]*GetResponse, 0)
	for _, video := range videos {
		response = append(response, NewGetResponse(video))
	}

	return response
}

type Quality struct {
	ID       uint            `json:"-"`
	Code     string          `json:"code"`
	Title    string          `json:"title"`
	Type     QualityType     `json:"-"`
	Settings QualitySettings `json:"-"`
}

type QualityType struct {
	Type string
}

var QualityTypeResize = QualityType{"Resize"}
var QualityTypeConvert = QualityType{"Convert"}
var QualityTypeSkip = QualityType{"Skip"}

func (q *Quality) Process(ctx context.Context, pathFrom, pathTo string) error {
	switch q.Type {
	case QualityTypeResize:
		audioBitrate := q.Settings.AudioBitrate
		if audioBitrate == "" {
			origAudioBitrate, _ := ffmpegthumbs.GetAudioBitrate(pathFrom)
			audioBitrate = origAudioBitrate
		}

		return ffmpegthumbs.Resize(
			ctx,
			q.Settings.Height,
			q.Settings.CRF,
			q.Settings.Speed,
			q.Settings.VideoBitrate,
			audioBitrate,
			pathFrom,
			pathTo,
			q.Code,
		)
	case QualityTypeConvert:
		origVideoBitrate, err := ffmpegthumbs.GetVideoBitrate(pathFrom)
		if err != nil {
			origVideoBitrate = "30000"
		}

		return ffmpegthumbs.ToWebm(
			ctx,
			pathFrom,
			q.Settings.CRF,
			q.Settings.Speed,
			origVideoBitrate,
			pathTo,
			q.Code,
		)
	}

	return nil
}

type QualitySettings struct {
	MinHeight    int
	Height       int
	CRF          string
	Speed        string
	VideoBitrate string
	AudioBitrate string
}

func GetQuality(id uint) *Quality {
	for _, quality := range SupportedQualities {
		if quality.ID == id {
			return &quality
		}
	}

	return nil
}

var SupportedQualities = []Quality{
	{
		ID:    1,
		Type:  QualityTypeSkip,
		Code:  "tmp",
		Title: "VIDEO_QUALITY_SOURCE",
	},
	{
		ID:    2,
		Type:  QualityTypeResize,
		Code:  "shakal",
		Title: "VIDEO_QUALITY_SHAKAL",
		Settings: QualitySettings{
			MinHeight: 0, Height: 240, CRF: "50", Speed: "5", VideoBitrate: "5", AudioBitrate: "2000",
		},
	},
	{
		ID:    3,
		Type:  QualityTypeResize,
		Code:  "360",
		Title: "VIDEO_QUALITY_360",
		Settings: QualitySettings{
			MinHeight: 360, Height: 360, CRF: "33", Speed: "3", VideoBitrate: "900k",
		},
	},
	{
		ID:    4,
		Type:  QualityTypeResize,
		Code:  "480",
		Title: "VIDEO_QUALITY_480",
		Settings: QualitySettings{
			MinHeight: 480, Height: 480, CRF: "33", Speed: "3", VideoBitrate: "1000k",
		},
	},
	{
		ID:    5,
		Type:  QualityTypeResize,
		Code:  "720",
		Title: "VIDEO_QUALITY_720",
		Settings: QualitySettings{
			MinHeight: 720, Height: 720, CRF: "32", Speed: "2", VideoBitrate: "1800k",
		},
	},
	{
		ID:    6,
		Type:  QualityTypeConvert,
		Code:  "origWebm",
		Title: "VIDEO_QUALITY_SOURCE",
		Settings: QualitySettings{
			MinHeight: 0, CRF: "31", Speed: "1",
		},
	},
}
