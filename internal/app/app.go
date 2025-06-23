package app

import (
	"context"
	"database/sql"
	"fmt"

	"log"
	"net/http"
	"os"

	"github.com/LuisDavid01/Websockets-go/internal/api"
	"github.com/LuisDavid01/Websockets-go/internal/manager"
	"github.com/LuisDavid01/Websockets-go/internal/store"

	"github.com/LuisDavid01/Websockets-go/migrations"
)

type Application struct {
	Logger  *log.Logger
	DB      *sql.DB
	Manager *manager.Manager
	Users   *api.UserHandler
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	ctx := context.Background()
	// our handler will go here

	//Database conection
	pgDb, err := store.Open()
	if err != nil {
		return nil, err
	}
	userStore := store.NewDBConn(pgDb)
	userHandler := api.NewUserHandler(userStore, logger)
	manager := manager.NewManager(ctx, userStore)

	//do migrations

	err = store.MigrateFS(pgDb, migrations.FS, ".")
	if err != nil {
		panic(err)
	}
	//we construct the application
	app := &Application{
		Logger:  logger,
		Manager: manager,
		DB:      pgDb,
		Users:   userHandler,
	}
	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is avaliable\n")
}
