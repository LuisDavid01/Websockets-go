package routes

import (
	//"io"
	//"html/template"
	//"net/http"

	"github.com/LuisDavid01/Websockets-go/internal/app"
	"github.com/LuisDavid01/Websockets-go/internal/templates"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()
	tmpl := templates.NewTemplates()
	appIndex := NewAppRoutes(tmpl)

	// Health check route
	r.Get("/health", app.HealthCheck)
	//route enableing ws
	r.Get("/ws", app.Manager.ServeWS)

	//auth
	r.Post("/login", app.Manager.LoginHandler)

	//serve frontend
	r.Get("/", appIndex.IndexHandler)
	return r
}
