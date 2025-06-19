package app

import (
	//"database/sql"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/LuisDavid01/Websockets-go/internal/manager"
)

type Application struct {
	Logger  *log.Logger
	Manager *manager.Manager
}

func NewApplication() (*Application, error) {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	ctx := context.Background()
	manager := manager.NewManager(ctx)
	// our handler will go here
	app := &Application{
		Logger:  logger,
		Manager: manager,
	}
	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is avaliable\n")
}
