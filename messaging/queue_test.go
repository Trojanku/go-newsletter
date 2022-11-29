package messaging_test

import (
	"Goo/integrationtest"
	"Goo/model"
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestQueue(t *testing.T) {
	integrationtest.SkipIfShort(t)

	t.Run("sends a message to the queue, receives it and delete it", func(t *testing.T) {

		queue, cleanup := integrationtest.CreateQueue()
		defer cleanup()

		err := queue.Send(context.Background(), model.Message{
			"foo": "bar",
		})
		require.NoError(t, err)

		m, receiptID, err := queue.Receive(context.Background())
		require.NoError(t, err)
		require.Equal(t, model.Message{"foo": "bar"}, *m)
		require.Greater(t, len(receiptID), 0)

		err = queue.Delete(context.Background(), receiptID)
		require.NoError(t, err)

		mr, _, err := queue.Receive(context.Background())
		require.NoError(t, err)
		require.Nil(t, mr)
	})
}
