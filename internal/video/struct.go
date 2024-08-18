package video

import (
	"database/sql/driver"
	"fmt"
	"golang.org/x/net/context"
	"gorm.io/gorm"
	"nine-dubz/internal/file"
	"nine-dubz/pkg/ffmpegthumbs"
	"reflect"
	"sort"
)

type Video struct {
	gorm.Model
	ID      uint
	Quality Quality `json:"-"`
	Title   string  `gorm:"-"`
	Width   int
	Height  int
	FileID  uint
	File    *file.File
}

type GetResponse struct {
	Quality Quality    `json:"quality"`
	Width   int        `json:"width"`
	Height  int        `json:"height"`
	File    *file.File `json:"file"`
}

func NewGetResponse(video Video) *GetResponse {
	return &GetResponse{
		Quality: video.Quality,
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

	sort.Slice(response, func(i, j int) bool {
		return response[i].Quality.Order < response[j].Quality.Order
	})

	return response
}

type Quality struct {
	ID       uint            `json:"-"`
	Code     string          `json:"code"`
	Title    string          `json:"title"`
	Type     QualityType     `json:"-"`
	Settings QualitySettings `json:"-"`
	Order    int             `json:"-"`
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

func (q *Quality) Scan(value interface{}) error {
	qualityId, ok := value.(int64)
	if !ok {
		return fmt.Errorf("can't scan category of type %v with id %v", reflect.TypeOf(value), value)
	}

	for _, quality := range SupportedQualities {
		if quality.ID == uint(qualityId) {
			*q = quality
			break
		}
	}

	return nil
}

func (q Quality) Value() (driver.Value, error) {
	for _, quality := range SupportedQualities {
		if q.ID == quality.ID {
			return int64(q.ID), nil
		}
	}

	return nil, fmt.Errorf("quality with id %d not found", q.ID)
}

var SupportedQualities = []Quality{
	{
		ID:    1,
		Type:  QualityTypeSkip,
		Code:  "tmp",
		Title: "VIDEO_QUALITY_SOURCE",
		Order: 1,
	},
	{
		ID:    2,
		Type:  QualityTypeResize,
		Code:  "shakal",
		Title: "VIDEO_QUALITY_SHAKAL",
		Order: 6,
		Settings: QualitySettings{
			MinHeight: 0, Height: 240, CRF: "50", Speed: "5", VideoBitrate: "5", AudioBitrate: "2000",
		},
	},
	{
		ID:    3,
		Type:  QualityTypeResize,
		Code:  "360",
		Title: "VIDEO_QUALITY_360",
		Order: 5,
		Settings: QualitySettings{
			MinHeight: 360, Height: 360, CRF: "33", Speed: "3", VideoBitrate: "900k",
		},
	},
	{
		ID:    4,
		Type:  QualityTypeResize,
		Code:  "480",
		Title: "VIDEO_QUALITY_480",
		Order: 4,
		Settings: QualitySettings{
			MinHeight: 480, Height: 480, CRF: "33", Speed: "3", VideoBitrate: "1000k",
		},
	},
	{
		ID:    5,
		Type:  QualityTypeResize,
		Code:  "720",
		Title: "VIDEO_QUALITY_720",
		Order: 3,
		Settings: QualitySettings{
			MinHeight: 720, Height: 720, CRF: "32", Speed: "2", VideoBitrate: "1800k",
		},
	},
	{
		ID:    6,
		Type:  QualityTypeConvert,
		Code:  "origWebm",
		Title: "VIDEO_QUALITY_SOURCE",
		Order: 2,
		Settings: QualitySettings{
			MinHeight: 0, CRF: "31", Speed: "1",
		},
	},
}
