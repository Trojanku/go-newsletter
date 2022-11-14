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
