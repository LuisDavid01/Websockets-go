package routes

import (
	//"io"
	"net/http"
	//"text/template"

	"github.com/LuisDavid01/Websockets-go/internal/app"
	"github.com/go-chi/chi/v5"
)

/*
	soon i may change to go templates

	type Templates struct {
		templates *template.Template
	}

	func (t *Templates) Render(w io.Writer, name string, data interface{}, c chi.Context) error {
		return t.templates.ExecuteTemplate(w, name, data)
	}

	func newTempalte() *Templates {
		return &Templates{
			templates: template.Must(template.ParseGlob("frontend/*.html")),
		}

}
*/
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
