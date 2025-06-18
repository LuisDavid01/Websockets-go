package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	websocketUpgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type Manager struct {
	clients ClientList
	sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		clients: make(ClientList),
	}
}

func (M *Manager) ServeWS(w http.ResponseWriter, r *http.Request) {
	log.Println("i have a new connection")
	//upgrade http to ws
	conn, err := websocketUpgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return
	}
	client := NewClient(conn, M)
	M.addClient(client)

	go client.readMessages()
	go client.writeMessages()

}

func (M *Manager) addClient(client *Client) {
	M.Lock()
	defer M.Unlock()
	M.clients[client] = true
}

func (M *Manager) removeClient(client *Client) {
	M.Lock()
	defer M.Unlock()
	if _, ok := M.clients[client]; ok {
		client.connection.Close()
		delete(M.clients, client)
	}
}
