package model

import (
	"fmt"
	"net/http"
)

type Pagination struct {
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
}

func (p *Pagination) Bind(r *http.Request) error {
	if p.Limit <= 0 {
		p.Limit = -1
	}

	if p.Offset <= 0 {
		p.Offset = -1
	}
	fmt.Println(p)
	return nil
}
