package store

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pressly/goose/v3"
)

func Open() (*sql.DB, error) {
	db, err := sql.Open("pgx", os.Getenv("POSTGRES_URL"))
	if err != nil {
		return nil, fmt.Errorf("Couldnt connecto to the database: %v", err)
	}

	fmt.Println("Database connected!")
	return db, nil
}

func MigrateFS(db *sql.DB, migrationsFS fs.FS, dir string) error {
	goose.SetBaseFS(migrationsFS)
	defer func() {
		goose.SetBaseFS(nil)
	}()

	return Migrate(db, dir)
}

func Migrate(db *sql.DB, dir string) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("Migration ERROR: %w", err)
	}

	err = goose.Up(db, dir)
	if err != nil {
		return fmt.Errorf("Migrating up ERROR: %w", err)
	}

	return nil
}
