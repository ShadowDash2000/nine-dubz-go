package usecase

import (
	"github.com/gorilla/websocket"
	"net/http"
)

type SocketRepository interface {
	UpgradeConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error)
	AddConnection(conn *websocket.Conn)
	DeleteConnection(conn *websocket.Conn)
	SendJson(i *interface{}, conn *websocket.Conn) error
}
