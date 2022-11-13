package server_test

import (
	"Goo/integrationtest"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

// Should be mocked
func TestServer_Start(t *testing.T) {
	// Failed: import cycle not allowed in test
	//server := Server{
	//	address: "test",
	//	mux:     nil,
	//	server:  nil,
	//}
	//t.Run("Should return error when server is nil", func(t *testing.T) {
	//	err := server.Start()
	//	require.NotNil(t, err)
	//})
	t.Run("Starts the server and listens for requests", func(t *testing.T) {
		integrationtest.SkipIfShort(t)

		cleanup := integrationtest.CreateServer()
		defer cleanup()

		resp, err := http.Get("http://localhost:8080/")
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
