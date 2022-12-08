package jobs_test

import (
	"Goo/jobs"
	"Goo/model"
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

type mockConfirmationEmailer struct {
	err   error
	to    model.Email
	token string
}

func (m *mockConfirmationEmailer) SendNewsletterConfirmationEmail(_ context.Context, to model.Email, token string) error {
	m.to = to
	m.token = token
	return m.err
}

func TestSendNewsletterConfirmationEmail(t *testing.T) {
	r := testRegistry{}

	t.Run("passes the recipient email and token to the email sender", func(t *testing.T) {

		emailer := &mockConfirmationEmailer{}
		jobs.SendNewsletterConfirmationEmail(r, emailer)

		job, ok := r["confirmation_email"]
		require.True(t, ok)

		err := job(context.Background(), model.Message{"email": "you@example.com", "token": "123"})
		require.NoError(t, err)

		require.Equal(t, "you@example.com", emailer.to.String())
		require.Equal(t, "123", emailer.token)
	})
}
