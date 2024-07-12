package model

type UploadStatus struct {
	Status int    `json:"status"`
	Error  string `json:"error,omitempty"`
}
