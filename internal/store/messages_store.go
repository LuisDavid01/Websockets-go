package store

import (
	"database/sql"
	"time"
)

type Chat struct {
	content []Messages
}
type Messages struct {
	ID             int64     `json:"id"`
	ReceptionistID int64     `json:"recept_id"`
	PacientID      int64     `json:"pacient_id"`
	Message        []byte    `json:"message_text"`
	IsRead         bool      `json:"is_read"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type IMessages interface {
	ReadChat(rID, pID int64) (*Chat, error)
	SendMessage(rID, PID int64) error
	BlockChat(chatID int64) error
	TransferChat(rID, chatID int64) error
	DownloadChat(chatID int64) ([]byte, error)
}

type DbConn struct {
	db *sql.DB
}

func NewDbConn(db *sql.DB) *DbConn {
	return &DbConn{db: db}
}
