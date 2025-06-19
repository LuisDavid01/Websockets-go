package store

import (
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type password struct {
	plainText string
	hash      []byte
}

// sets a encripted password
func (p *password) Set(plainText string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainText), 13)
	if err != nil {
		return err
	}
	p.plainText = plainText
	p.hash = hash
	return nil
}

// checks if a password is correct or nah
func (p *password) Matches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil // Password does not match
		default:
			return false, err //something went wrong in the server

		}
	}
	return true, nil
}

type User struct {
	ID        uint8     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Rol       string    `json:"rol"`
	Password  password  `json:"_"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// set a default anonimous user to compare to
var AnonUser = &User{}

func isAnon(u *User) bool {
	return u == AnonUser
}

type dbConn struct {
	db *sql.DB
}

func NewDBConn(db *sql.DB) *dbConn {
	return &dbConn{db: db}
}

// interface implement this later
type IUser interface {
	CreateUser(user *User) error
	GetUserById(username string) (*User, error)
	GetUserByUsername(id uint8) (*User, error)
	UpdateUser(*User) error
	GetUserToken(scope, tokenPlainText string) (*User, error)
}
