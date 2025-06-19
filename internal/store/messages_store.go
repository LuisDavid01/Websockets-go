package store

import "time"

type Chat struct {
	content []Messages
}
type Messages struct {
	ID             uint8     `json:"id"`
	ReceptionistID uint8     `json:"recept_id"`
	PacientID      uint8     `json:"pacient_id"`
	Message        []byte    `json:"message_text"`
	IsRead         bool      `json:"is_read"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type IMessage interface {
	ReadChat(rID, pID uint8) (*Chat, error)
	SendMessage(rID, PID uint8) error
	BlockChat(chatID uint8) error
	TransferChat(rID, chatID uint8) error
	DownloadChat(chatID uint8) ([]byte, error)
}
