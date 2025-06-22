package templates

import (
	"html/template"
	"net/http"
)

type Templates struct {
	Tmpl *template.Template
}

func NewTemplates() *Templates {
	return &Templates{
		Tmpl: template.Must(template.ParseGlob("frontend/*.html")),
	}
}

func (t *Templates) Render(w http.ResponseWriter, name string, data any) {
	err := t.Tmpl.ExecuteTemplate(w, name, data)
	if err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
	}
}
