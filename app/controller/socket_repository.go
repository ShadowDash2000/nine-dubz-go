package controller

import (
	"github.com/gorilla/websocket"
	"net/http"
)

type SocketRepository struct {
	Clients map[*websocket.Conn]bool
}

func (sr *SocketRepository) UpgradeConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upgrader := websocket.Upgrader{}
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	return connection, nil
}

func (sr *SocketRepository) AddConnection(conn *websocket.Conn) {
	sr.Clients[conn] = true
}

func (sr *SocketRepository) DeleteConnection(conn *websocket.Conn) {
	delete(sr.Clients, conn)
}

func (sr *SocketRepository) SendJson(i *interface{}, conn *websocket.Conn) error {
	return conn.WriteJSON(i)
}
