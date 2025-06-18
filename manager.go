package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type Manager struct {
}

func NewManager() *Manager {
	return &Manager{}
}

func (M *Manager) ServeWS(w http.ResponseWriter, r *http.Request) {
	log.Println("i have a new connection")
	//upgrade http to ws
	conn, err := websocketUpgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}

	conn.Close()
}
