package jobs

import (
	"Goo/model"
	"context"
	"errors"
	"fmt"
	"time"
)

type newsletterConfirmationEmailSender interface {
	SendNewsletterConfirmationEmail(ctx context.Context, to model.Email, token string) error
}

func SendNewsletterConfirmationEmail(r registry, es newsletterConfirmationEmailSender) {
	r.Register("confirmation_email", func(ctx context.Context, message model.Message) error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		to, ok := message["email"]
		if !ok {
			return errors.New("no email address in message")
		}

		token, ok := message["token"]
		if !ok {
			return errors.New("no token in message")
		}

		if err := es.SendNewsletterConfirmationEmail(ctx, model.Email(to), token); err != nil {
			return fmt.Errorf("error sending newsletter confirmation email: %w", err)
		}

		return nil
	})
}
