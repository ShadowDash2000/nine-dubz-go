package usecase

import (
	"github.com/gorilla/websocket"
	"net/http"
)

type SocketInteractor struct {
	SocketRepository SocketRepository
}

func (si *SocketInteractor) UpgradeConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	return si.SocketRepository.UpgradeConnection(w, r)
}

func (si *SocketInteractor) AddConnection(conn *websocket.Conn) {
	si.SocketRepository.AddConnection(conn)
}

func (si *SocketInteractor) DeleteConnection(conn *websocket.Conn) {
	si.SocketRepository.DeleteConnection(conn)
}

func (si *SocketInteractor) SendJson(i *interface{}, conn *websocket.Conn) error {
	return si.SocketRepository.SendJson(i, conn)
}
