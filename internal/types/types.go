package types

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
)

// ClientList is a map of clients.
type ClientList map[*Client]bool

// Manager interface defines methods that the Client needs from the Manager.
type Manager interface {
	RemoveClient(*Client)
	RouteEvent(Event, *Client) error
}

// Client represents a WebSocket client.
type Client struct {
	Connection *websocket.Conn
	Manager    Manager
	Chatroom   string
	Egress     chan Event
}

// Event represents a WebSocket event.
type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

// EventHandler defines the signature for handling events.
type EventHandler func(event Event, c *Client) error

// Event types.
const (
	EventSendMessage = "send_message"
	EventNewMessage  = "new_message"
	EventChatRoom    = "change_room"
)

// SendMessageEvent represents a message sent by a client.
type SendMessageEvent struct {
	Message string `json:"message"`
	From    string `json:"from"`
}

// NewMessageEvent represents a broadcasted message.
type NewMessageEvent struct {
	SendMessageEvent
	Sent time.Time `json:"sent"`
}

// ChangeRoomEvent represents a request to change chatrooms.
type ChangeRoomEvent struct {
	Name string `json:"name"`
}
