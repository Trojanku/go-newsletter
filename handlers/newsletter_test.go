package handlers_test

import (
	"Goo/handlers"
	"Goo/model"
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type signupperMock struct {
	email model.Email
}

func (s *signupperMock) SignupForNewsletter(ctx context.Context, email model.Email) (string, error) {
	s.email = email
	return "", nil
}

func TestNewsletterSignup(t *testing.T) {
	mux := chi.NewMux()
	s := &signupperMock{}
	handlers.NewsletterSignup(mux, s)

	t.Run("signs up a valid email address", func(t *testing.T) {
		code, _, _ := makePostRequest(mux, "/newsletter/signup", createFormHeader(),
			strings.NewReader("email=me%40example.com"))
		require.Equal(t, http.StatusFound, code)
		require.Equal(t, model.Email("me@example.com"), s.email)
	})

	t.Run("rejects an invalid email address", func(t *testing.T) {
		code, _, _ := makePostRequest(mux, "/newsletter/signup", createFormHeader(),
			strings.NewReader("email=notanemail"))
		require.Equal(t, http.StatusBadRequest, code)
	})
}

func makePostRequest(handler http.Handler, target string, header http.Header, body io.Reader) (int, http.Header, string) {
	req := httptest.NewRequest(http.MethodPost, target, body)
	req.Header = header
	res := httptest.NewRecorder()
	handler.ServeHTTP(res, req)
	result := res.Result()
	bodyBytes, err := io.ReadAll(result.Body)
	if err != nil {
		panic(err)
	}
	return result.StatusCode, result.Header, string(bodyBytes)
}

func createFormHeader() http.Header {
	header := http.Header{}
	header.Set("Content-Type", "application/x-www-form-urlencoded")
	return header
}
