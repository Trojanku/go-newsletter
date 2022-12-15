package storage_test

import (
	"Goo/integrationtest"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDatabase_SignupForNewsletter(t *testing.T) {
	integrationtest.SkipIfShort(t)

	t.Run("signs up", func(t *testing.T) {
		db, cleanup := integrationtest.CreateDatabase()
		defer cleanup()

		expectedToken, err := db.SignupForNewsletter(context.Background(), "me@example.com")
		require.NoError(t, err)
		require.Equal(t, 64, len(expectedToken))

		var email, token string
		err = db.DB.QueryRow(`select email, token from newsletter_subscribers `).Scan(&email, &token)
		require.NoError(t, err)
		require.Equal(t, "me@example.com", email)
		assert.Equal(t, expectedToken, token)

		expectedToken2, err := db.SignupForNewsletter(context.Background(), "me@example.com")
		require.NoError(t, err)
		require.NotEqual(t, expectedToken, expectedToken2)

		err = db.DB.QueryRow(`select email, token from newsletter_subscribers`).Scan(&email, &token)
		require.NoError(t, err)
		require.Equal(t, "me@example.com", email)
		assert.Equal(t, expectedToken2, token)
	})
}

func TestDatabase_ConfirmNewsletterSignup(t *testing.T) {
	integrationtest.SkipIfShort(t)

	t.Run("confirms subscriber from the token and returns the associated email address", func(t *testing.T) {
		db, cleanup := integrationtest.CreateDatabase()
		defer cleanup()

		token, err := db.SignupForNewsletter(context.Background(), "me@example.com")
		require.NoError(t, err)

		var confirmed bool
		err = db.DB.Get(&confirmed, `select confirmed from newsletter_subscribers where token = &1`, token)
		require.NoError(t, err)
		require.False(t, confirmed)

		email, err := db.ConfirmNewsletterSignup(context.Background(), token)
		require.NoError(t, err)
		require.Equal(t, "me@example.com", email.String())

		err = db.DB.Get(&confirmed, `select confirmed from newsletter_subscribers where token = &1`, token)
		require.NoError(t, err)
		require.True(t, confirmed)
	})

	t.Run("returns nil if no such token", func(t *testing.T) {
		db, cleanup := integrationtest.CreateDatabase()
		defer cleanup()

		_, err := db.SignupForNewsletter(context.Background(), "me@example.com")
		require.NoError(t, err)

		email, err := db.ConfirmNewsletterSignup(context.Background(), "wrongtoken")
		require.NoError(t, err)
		require.Nil(t, email)
	})
}
