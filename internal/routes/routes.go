package routes

import (
	"net/http"

	"github.com/LuisDavid01/Websockets-go/internal/app"
	"github.com/go-chi/chi/v5"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()
	// Health check route
	r.Get("/health", app.HealthCheck)
	//route enableing ws
	r.Get("/ws", app.Manager.ServeWS)

	//auth
	r.Post("/login", app.Manager.LoginHandler)

	//serve frontend
	fs := http.FileServer(http.Dir("./frontend"))
	r.Handle("/*", http.StripPrefix("/", fs))
	return r
}
