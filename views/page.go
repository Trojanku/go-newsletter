package views

import (
	"Goo/templates"
	"html/template"
)

func LoadTemplate() (*template.Template, error) {
	tmpl, err := template.New("").Parse(templates.Index)
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}
