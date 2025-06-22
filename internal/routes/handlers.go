package routes

import (
	"net/http"

	"github.com/LuisDavid01/Websockets-go/internal/templates"
	"github.com/LuisDavid01/Websockets-go/internal/types"
)

type AppRoutes struct {
	Templates *templates.Templates
}

func NewAppRoutes(t *templates.Templates) *AppRoutes {
	return &AppRoutes{Templates: t}
}

func (a *AppRoutes) IndexHandler(w http.ResponseWriter, r *http.Request) {
	page := types.Page{
		Title: "Inicio",
		User:  "Invitado",
		Form:  types.NewFormData(),
	}
	a.Templates.Render(w, "index.html", page)
}
