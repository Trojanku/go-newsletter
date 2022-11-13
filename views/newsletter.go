package views

import (
	"Goo/templates"
	"html/template"
)

func NewsletterThanksPage(path string) (*template.Template, error) {
	return template.New("").Parse(templates.Thanks)
}
