package video

import (
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
	Title  *string    `json:"title"`
	Width  int        `json:"width"`
	Height int        `json:"height"`
	File   *file.File `json:"file"`
}

func NewGetResponse(video Video) *GetResponse {
	return &GetResponse{
		Title:  GetQualityTitle(video.QualityID),
		Width:  video.Width,
		Height: video.Height,
		File:   video.File,
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
	ID       uint
	Title    string
	Type     QualityType
	Settings QualitySettings
}

type QualityType struct {
	Type string
}

var QualityTypeResize = QualityType{"Resize"}
var QualityTypeConvert = QualityType{"Convert"}
var QualityTypeSkip = QualityType{"Skip"}

func (q *Quality) Process(pathFrom, pathTo string) error {
	switch q.Type {
	case QualityTypeResize:
		audioBitrate := q.Settings.AudioBitrate
		if audioBitrate == "" {
			origAudioBitrate, _ := ffmpegthumbs.GetAudioBitrate(pathFrom)
			audioBitrate = origAudioBitrate
		}

		return ffmpegthumbs.Resize(
			q.Settings.Height,
			q.Settings.CRF,
			q.Settings.Speed,
			q.Settings.VideoBitrate,
			audioBitrate,
			pathFrom,
			pathTo,
			q.Title,
		)
	case QualityTypeConvert:
		origVideoBitrate, err := ffmpegthumbs.GetVideoBitrate(pathFrom)
		if err != nil {
			origVideoBitrate = "30000"
		}

		return ffmpegthumbs.ToWebm(
			pathFrom,
			q.Settings.CRF,
			q.Settings.Speed,
			origVideoBitrate,
			pathTo,
			q.Title,
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

func GetQualityTitle(id uint) *string {
	for _, quality := range SupportedQualities {
		if quality.ID == id {
			return &quality.Title
		}
	}

	return nil
}

var SupportedQualities = []Quality{
	{
		ID:    1,
		Type:  QualityTypeSkip,
		Title: "tmp",
	},
	{
		ID:    2,
		Type:  QualityTypeResize,
		Title: "shakal",
		Settings: QualitySettings{
			MinHeight: 0, Height: 240, CRF: "50", Speed: "5", VideoBitrate: "5", AudioBitrate: "2000",
		},
	},
	{
		ID:    3,
		Type:  QualityTypeResize,
		Title: "360",
		Settings: QualitySettings{
			MinHeight: 360, Height: 360, CRF: "33", Speed: "3", VideoBitrate: "900k",
		},
	},
	{
		ID:    4,
		Type:  QualityTypeResize,
		Title: "480",
		Settings: QualitySettings{
			MinHeight: 480, Height: 480, CRF: "33", Speed: "3", VideoBitrate: "1000k",
		},
	},
	{
		ID:    5,
		Type:  QualityTypeResize,
		Title: "720",
		Settings: QualitySettings{
			MinHeight: 720, Height: 720, CRF: "32", Speed: "2", VideoBitrate: "1800k",
		},
	},
	{
		ID:    6,
		Type:  QualityTypeConvert,
		Title: "origWebm",
		Settings: QualitySettings{
			MinHeight: 0, CRF: "31", Speed: "1",
		},
	},
}
