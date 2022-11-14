package handlers_test

import (
	"Goo/handlers"
	"context"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type pingerMock struct {
	err error
}

func (p *pingerMock) Ping(_ context.Context) error {
	return p.err
}

func TestHealth(t *testing.T) {
	t.Run("returns 200", func(t *testing.T) {
		mux := chi.NewMux()

		handlers.Health(mux, &pingerMock{})
		code, _, _ := makeGetRequest(mux, "/health")
		require.Equal(t, http.StatusOK, code)
	})
	t.Run("returns 502 if the database cannot be pinged", func(t *testing.T) {
		mux := chi.NewMux()
		handlers.Health(mux, &pingerMock{err: errors.New("error")})
		code, _, _ := makeGetRequest(mux, "/health")
		require.Equal(t, http.StatusBadGateway, code)
	})
}

func makeGetRequest(handler http.Handler, target string) (int, http.Header, string) {
	req := httptest.NewRequest(http.MethodGet, target, nil)
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)
	result := res.Result()
	bodyBytes, err := io.ReadAll(result.Body)
	if err != nil {
		panic(err)
	}
	return result.StatusCode, result.Header, string(bodyBytes)
}
