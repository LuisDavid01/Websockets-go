package client

import (
	"encoding/json"
	"log"
	"time"

	"github.com/LuisDavid01/Websockets-go/internal/types"
	"github.com/gorilla/websocket"
)

var (
	pongWait     = 10 * time.Second
	pingInterval = (pongWait * 9) / 10
)

// NewClient crea un nuevo cliente.
func NewClient(conn *websocket.Conn, manager types.Manager) *types.Client {
	return &types.Client{
		Connection: conn,
		Manager:    manager,
		Egress:     make(chan types.Event),
	}
}

// ReadMessages lee mensajes de la conexión WebSocket.
func ReadMessages(c *types.Client) {
	defer func() {
		c.Manager.RemoveClient(c)
	}()

	if err := c.Connection.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		log.Println("Timeout err", err)
		return
	}
	c.Connection.SetReadLimit(512)
	c.Connection.SetPongHandler(func(pongMsg string) error {
		return PongHandler(c, pongMsg)
	})

	for {
		_, payload, err := c.Connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("Connection closed", err)
			}
			break
		}

		var request types.Event
		if err := json.Unmarshal(payload, &request); err != nil {
			log.Printf("ERROR: %v", err)
			break
		}

		if err := c.Manager.RouteEvent(request, c); err != nil {
			log.Println("Error handling the request:", err)
		}
	}
}

// WriteMessages escribe mensajes en la conexión WebSocket.
func WriteMessages(c *types.Client) {
	defer func() {
		c.Manager.RemoveClient(c)
	}()

	ticker := time.NewTicker(pingInterval)
	defer ticker.Stop()

	for {
		select {
		case message, ok := <-c.Egress:
			if !ok {
				if err := c.Connection.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("Connection is closed", err)
				}
				return
			}
			data, err := json.Marshal(message)
			if err != nil {
				log.Printf("ERROR parsing the json data: %v", err)
				return
			}
			if err := c.Connection.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("Failed to send message: %v", err)
			}
			log.Println("Message sent")
		case <-ticker.C:
			log.Println("ping")
			if err := c.Connection.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				log.Println("Writing ping error: ", err)
				return
			}
		}
	}
}

// PongHandler maneja mensajes pong.
func PongHandler(c *types.Client, pongMsg string) error {
	log.Println("I received a pong")
	return c.Connection.SetReadDeadline(time.Now().Add(pongWait))
}
