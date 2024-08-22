package etag

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"hash"
	"net/http"
	"strconv"
)

type HashWriter struct {
	w      http.ResponseWriter
	hash   hash.Hash
	buff   *bytes.Buffer
	len    int
	status int
}

func (hw *HashWriter) Header() http.Header {
	return hw.w.Header()
}

func (hw *HashWriter) WriteHeader(status int) {
	hw.status = status
}

func (hw *HashWriter) Write(b []byte) (int, error) {
	if hw.status == 0 {
		hw.status = http.StatusOK
	}

	hw.buff.Write(b)

	l, err := hw.hash.Write(b)
	hw.len += l
	return l, err
}

func Etag(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Upgrade") != "" {
			next.ServeHTTP(w, r)
			return
		}

		hw := HashWriter{
			w:    w,
			hash: sha1.New(),
			buff: bytes.NewBuffer(nil),
		}

		next.ServeHTTP(&hw, r)

		if strconv.Itoa(hw.status)[0] != '2' || hw.buff.Len() == 0 {
			if hw.status > 0 {
				w.WriteHeader(hw.status)
			}
			w.Write(hw.buff.Bytes())
			return
		}

		etag := fmt.Sprintf("W/%v-%v", strconv.Itoa(hw.len), hex.EncodeToString(hw.hash.Sum(nil)))

		w.Header().Set("ETag", etag)

		if r.Header.Get("If-None-Match") == etag {
			w.WriteHeader(http.StatusNotModified)
			return
		} else if hw.status > 0 {
			w.WriteHeader(hw.status)
		}
		w.Write(hw.buff.Bytes())
	})
}
