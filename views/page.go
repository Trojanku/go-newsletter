package views

import "html/template"

func LoadTemplate(path string) (*template.Template, error) {
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}
