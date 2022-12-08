package messaging

import (
	"Goo/model"
	"context"
	"embed"
	"fmt"
	"go.uber.org/zap"
	"strings"
)
import "github.com/go-gomail/gomail"

//go:embed emails
var emails embed.FS

type Emailer struct {
	baseURL           string
	marketingFrom     string
	transactionalFrom string

	dialer *gomail.Dialer
	log    *zap.Logger
}

type NewEmailerOptions struct {
	BaseURL string

	Host string
	Port int

	MarketingUsername string
	MarketingPassword string

	TransactionalUsername string
	TransactionalPassword string

	MarketingEmailAddress string
	MarketingEmailName    string

	TransactionalEmailAddress string
	TransactionalEmailName    string

	Log *zap.Logger
}

func NewEmailer(opts NewEmailerOptions) *Emailer {
	return &Emailer{
		baseURL: opts.BaseURL,

		marketingFrom:     opts.MarketingEmailName,
		transactionalFrom: opts.TransactionalEmailName,

		dialer: setupDialer(dialerSettings{
			Host:     opts.Host,
			Port:     opts.Port,
			Username: opts.TransactionalUsername,
			Password: opts.TransactionalPassword,
		}),
		log: opts.Log,
	}
}

type dialerSettings struct {
	Host         string
	Port         int
	Username     string
	Password     string
	WriteTimeout string
}

func setupDialer(settings dialerSettings) *gomail.Dialer {
	return gomail.NewDialer(settings.Host, settings.Port, settings.Username, settings.Password)
}

// SendNewsletterConfirmationEmail with a confirmation link.
// This is a transactional email, because it's a response to a user action.
func (e *Emailer) SendNewsletterConfirmationEmail(_ context.Context, to model.Email, token string) error {
	keywords := map[string]string{
		"base_url":   e.baseURL,
		"action_url": e.baseURL + "/newsletter/confirm?token=" + token,
	}

	return e.send(requestBody{
		From:      e.transactionalFrom,
		ToAddress: to.String(),
		// TODO: change to name
		ToName:      to.String(),
		Subject:     "Confirm your subscription to the newsletter",
		ContentHTML: getEmail("confirmation_email.html", keywords),
		ContextText: getEmail("confirmation_email.txt", keywords),
	})
}

type requestBody struct {
	From        string
	ToAddress   string
	ToName      string
	Subject     string
	ContentHTML string
	ContextText string
}

func (e *Emailer) send(body requestBody) error {
	m := gomail.NewMessage()

	m.SetHeader("From", body.From)
	m.SetAddressHeader("To", body.ToAddress, body.ToName)
	m.SetHeader("Subject", body.Subject)
	m.SetBody("text/html", body.ContentHTML)
	m.SetBody("text/plain", body.ContextText)

	if err := e.dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("could not send email: %w", err)
	}
	return nil
}

// getEmail from the given path, panicking on errors.
// It also replaces keywords given in the map.
func getEmail(path string, keywords map[string]string) string {
	email, err := emails.ReadFile("emails/" + path)
	if err != nil {
		panic(err)
	}

	emailString := string(email)
	for keyword, replacement := range keywords {
		emailString = strings.ReplaceAll(emailString, "{{"+keyword+"}}", replacement)
	}

	return emailString
}
