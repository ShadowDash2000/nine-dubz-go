package model

import "gorm.io/gorm"

type ApiMethod struct {
	gorm.Model
	Path   string `json:"path"`
	Method string `json:"method"`
}
