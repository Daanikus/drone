package main

import (
	"html/template"
	"io"

	"github.com/labstack/echo"
)

type Template struct {
	templates *template.Template
}

func buildRenderer() (*Template, error) {
	tmpl, err := template.ParseGlob("tmpl/*.tmpl")
	if err != nil {
		return nil, err
	}
	t := &Template{
		templates: tmpl,
	}
	return t, nil
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
