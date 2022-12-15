package views

import (
	"Goo/templates"
	"html/template"
)

func NewsletterThanksPage(path string) (*template.Template, error) {
	return template.New(path).Parse(templates.Thanks)
}

func NewsletterConfirmPage(path string) (*template.Template, error) {
	return template.New(path).Parse(templates.Confirm)
}

func NewsletterConfirmedPage(path string) (*template.Template, error) {
	return template.New(path).Parse(templates.Confirmed)
}
