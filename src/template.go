package main

import (
	"html/template"
	"io"

	"github.com/labstack/echo"
)

type TemplateRenderer struct { //nolint: unused
	templates       *template.Template
	templatePath    string
	reloadTemplates bool
}

func NewTemplateRenderer(templatePath string, reloadTemplates bool) *TemplateRenderer { //nolint
	return &TemplateRenderer{
		templatePath:    templatePath,
		reloadTemplates: reloadTemplates,
		templates:       template.Must(template.ParseGlob(templatePath)),
	}
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	if t.reloadTemplates {
		t.templates = template.Must(template.ParseGlob(t.templatePath))
	}

	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}
