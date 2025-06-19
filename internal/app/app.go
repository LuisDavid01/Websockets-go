package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/LuisDavid01/Websockets-go/internal/manager"
	"github.com/LuisDavid01/Websockets-go/internal/store"
)

type Application struct {
	Logger  *log.Logger
	DB      *sql.DB
	Manager *manager.Manager
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	ctx := context.Background()
	manager := manager.NewManager(ctx)
	// our handler will go here

	//Database conection
	pgDb, err := store.Open()
	if err != nil {
		return nil, err
	}

	//do migrations

	//we construct the application
	app := &Application{
		Logger:  logger,
		Manager: manager,
		DB:      pgDb,
	}
	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is avaliable\n")
}
