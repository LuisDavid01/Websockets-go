package event

import (
	"encoding/json"
	"time"

	"github.com/LuisDavid01/Websockets-go/internal/types"
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(event Event, c *types.Client) error

const (
	EventSendMessage = "send_message"
	EventNewMessage  = "new_message"
	EventChatRoom    = "change_room"
)

type SendMessageEvent struct {
	Message string `json:"message"`
	From    string `json:"from"`
}

type NewMessageEvent struct {
	SendMessageEvent
	Sent time.Time `json:"sent"`
}

type ChangeRoomEvent struct {
	Name string `json:"name"`
}
