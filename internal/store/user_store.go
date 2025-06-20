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
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Rol       string    `json:"-"`
	Password  password  `json:"_"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// set a default anonimous user to compare to
var AnonUser = &User{}

func isAnon(user *User) bool {
	return user == AnonUser
}

type dbConn struct {
	db *sql.DB
}

func NewDBConn(db *sql.DB) *dbConn {
	return &dbConn{db: db}
}

// interface implement this later
type IUser interface {
	RegisterUser(user *User) error
	GetUserByUsername(username string) (*User, error)
	GetUserById(id int64) (*User, error)
	UpdateUser(user *User) error
	GetUserToken(scope, tokenPlainText string) (*User, error)
}

func (pg *dbConn) RegisterUser(user *User) error {
	query := ` INSERT INTO 
	users(username, email, rol, password)
	VALUES($1,$2,$3,$4)
	`
	err := pg.db.QueryRow(query, user.Username, user.Email, user.Rol, user.Password.hash).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return err
	}
	return nil
}

func (pg *dbConn) GetUserById(id int64) (*User, error) {

	user := User{
		Password: password{},
	}
	query := `SELECT id, username, password, email, "user", created_at, updated_at
	FROM users
	WHERE id = $1
	`

	err := pg.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Password.hash,
		&user.Email,
		&user.Rol,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	//not found
	if err == sql.ErrNoRows {
		return nil, nil
	}
	// oops
	if err != nil {
		return nil, err
	}
	//we found the user
	return &user, nil
}

func (pg *dbConn) GetUserByUsername(username string) (*User, error) {

	user := User{
		Password: password{},
	}
	query := `SELECT id, username, password, email, rol, created_at, updated_at
	FROM users
	WHERE username = $1
	`

	err := pg.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password.hash,
		&user.Email,
		&user.Rol,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	//not found
	if err == sql.ErrNoRows {
		return nil, nil
	}
	// oops
	if err != nil {
		return nil, err
	}
	//we found the user
	return &user, nil
}

func (pg *dbConn) UpdateUser(user *User) error {
	query := `
	UPDATE users
	SET username = $1, password = $2, email = $3, updated_at = CURRENT_TIMESTAMP
	WHERE id = $4
	returning updated_at
	`
	result, err := pg.db.Exec(query,
		user.Username,
		user.Password.hash,
		user.Email,
		user.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}
