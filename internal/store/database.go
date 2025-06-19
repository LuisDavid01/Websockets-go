package store

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
)

func Open() (*sql.DB, error) {
	db, err := sql.Open("pgx", os.Getenv("POSTGRES_URL"))
	if err != nil {
		return nil, fmt.Errorf("Couldnt connecto to the database: %v", err)
	}

	fmt.Println("Database connected!")
	return db, nil
}
