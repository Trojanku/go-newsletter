package handlers_test

import (
	"Goo/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"net/http"
	"strings"
	"testing"
)

func TestAddMetrics(t *testing.T) {
	t.Run("adds counter and histogram metrics", func(t *testing.T) {
		mux := chi.NewMux()
		registry := prometheus.NewRegistry()
		mux.Use(handlers.AddMetrics(registry))
		handlers.Metrics(mux, registry)
		mux.Get("/exists", func(writer http.ResponseWriter, request *http.Request) {

			code, _, _ := makeGetRequest(mux, "/exists")
			assert.Equal(t, http.StatusOK, code)
			code, _, _ = makeGetRequest(mux, "/doesnotexist")
			assert.Equal(t, http.StatusNotFound, code)

			code, _, body := makeGetRequest(mux, "/metrics")
			assert.Equal(t, http.StatusOK, code)

			assert.True(t, strings.Contains(body, `app_http_requests_total{code="200",method="GET", path="/exists"} 1`))
			assert.True(t, strings.Contains(body, `app_http_requests_total{code="404",method="GET", path="/doesnotexist"} 1`))

			assert.True(t, strings.Contains(body, `app_http_request_duration_seconds_bucket{code="200",le="+Inf"} 1`))
			assert.True(t, strings.Contains(body, `app_http_request_duration_seconds_bucket{code="404",le="+Inf"} 1`))
		})
	})
}