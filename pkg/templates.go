package pkg

import (
	"html/template"
	"io"
	"io/fs"
)

var templateFS fs.FS

func InitTemplates(fs fs.FS) {
	templateFS = fs
}

type tmpl struct {
	Title   string
	Payload any
}

func RenderTemplate(w io.Writer, page, title string, payload any) error {
	templ := template.Must(template.ParseFS(
		templateFS,
		"templates/pages/"+page+".tmpl",
		"templates/top.tmpl",
		"templates/bottom.tmpl",
	))
	return templ.Execute(w, tmpl{
		Title:   title,
		Payload: payload,
	})
}
