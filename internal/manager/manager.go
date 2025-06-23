package manager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/LuisDavid01/Websockets-go/internal/auth"
	"github.com/LuisDavid01/Websockets-go/internal/client"
	"github.com/LuisDavid01/Websockets-go/internal/store"
	"github.com/LuisDavid01/Websockets-go/internal/types"
	"github.com/gorilla/websocket"
)

var (
	websocketUpgrader = websocket.Upgrader{
		CheckOrigin:     checkOrigin,
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

// Manager manages WebSocket clients and events.
type Manager struct {
	clients types.ClientList
	sync.RWMutex

	otps      auth.RetentionMap
	handlers  map[string]types.EventHandler
	templates *template.Template
	store     store.IUser
}

// NewManager creates a new Manager.
func NewManager(ctx context.Context, store store.IUser) *Manager {
	m := &Manager{
		clients:  make(types.ClientList),
		handlers: make(map[string]types.EventHandler),
		otps:     auth.NewRetentionMap(ctx, 5*time.Second),
		store:    store,
	}
	m.setupEventHandlers()
	return m
}

func (m *Manager) setupEventHandlers() {
	m.handlers[types.EventSendMessage] = SendMessage
	m.handlers[types.EventChatRoom] = ChatRoomHandler
}

// ChatRoomHandler handles chat room change events.
func ChatRoomHandler(event types.Event, c *types.Client) error {
	var changeRoomEvent types.ChangeRoomEvent
	if err := json.Unmarshal(event.Payload, &changeRoomEvent); err != nil {
		return fmt.Errorf("bad payload: %v", err)
	}
	fmt.Printf("this is the chatroom: %v\n", changeRoomEvent.Name)
	c.Chatroom = changeRoomEvent.Name
	return nil
}

// SendMessage handles sending messages to clients in the same chatroom.
func SendMessage(event types.Event, c *types.Client) error {
	var chatEvent types.SendMessageEvent
	if err := json.Unmarshal(event.Payload, &chatEvent); err != nil {
		return fmt.Errorf("bad payload: %v", err)
	}

	var broadMessage types.NewMessageEvent
	broadMessage.Sent = time.Now()
	broadMessage.Message = chatEvent.Message
	broadMessage.From = chatEvent.From
	data, err := json.Marshal(broadMessage)
	if err != nil {
		return fmt.Errorf("failed to broadcast the message: %v", err)
	}

	outgoingEvent := types.Event{
		Payload: data,
		Type:    types.EventNewMessage,
	}

	for client := range c.Manager.(*Manager).clients {
		if client.Chatroom == c.Chatroom {
			client.Egress <- outgoingEvent
		}
	}
	return nil
}

// RouteEvent routes an event to the appropriate handler.
func (m *Manager) RouteEvent(event types.Event, c *types.Client) error {
	if handler, ok := m.handlers[event.Type]; ok {
		if err := handler(event, c); err != nil {
			return err
		}
		return nil
	}
	return errors.New("there is no such event")
}

// ServeWS handles WebSocket connections.
func (m *Manager) ServeWS(w http.ResponseWriter, r *http.Request) {
	log.Println("I have a new connection")

	otp := r.URL.Query().Get("otp")
	if otp == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if !m.otps.VerifyOTP(otp) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	conn, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}

	newClient := client.NewClient(conn, m)
	m.addClient(newClient)

	go client.ReadMessages(newClient)
	go client.WriteMessages(newClient)
}

// LoginHandler handles user login and OTP generation.
func (m *Manager) LoginHandler(w http.ResponseWriter, r *http.Request) {

	type userLoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var req userLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	lookupUser, err := m.store.GetUserByUsername(req.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if lookupUser == nil {
		http.Error(w, "Missing user or password", http.StatusInternalServerError)
		return

	}

	passwordDoMatch, err := lookupUser.Password.Matches(req.Password)

	if err != nil {
		http.Error(w, "error matching the password", http.StatusInternalServerError)
		return

	}
	if !passwordDoMatch {
		http.Error(w, "could not validate the user or password", http.StatusUnauthorized)
		return

	}

	type response struct {
		Username string `json:"username"`
		OTP      string `json:"otp"`
	}
	otp := m.otps.NewOTP()
	resp := response{
		OTP:      otp.Key,
		Username: req.Username,
	}
	data, err := json.Marshal(resp)
	if err != nil {
		log.Println("Failed to marshal OTP response:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
	return

}

// addClient adds a client to the manager.
func (m *Manager) addClient(client *types.Client) {
	m.Lock()
	defer m.Unlock()
	m.clients[client] = true
}

// RemoveClient removes a client from the manager.
func (m *Manager) RemoveClient(client *types.Client) {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.clients[client]; ok {
		client.Connection.Close()
		delete(m.clients, client)
	}
}

// checkOrigin validates the request origin.
func checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	switch origin {
	case "https://localhost:8080":
		return true
	default:
		return false
	}
}
