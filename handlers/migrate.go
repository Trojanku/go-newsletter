package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type migrator interface {
	MigrateTo(ctx context.Context, version uint) error
	MigrateUp(ctx context.Context) error
}

func MigrateTo(mux chi.Router, m migrator) {
	mux.Post("/migrate/to", func(w http.ResponseWriter, r *http.Request) {
		version := r.FormValue("version")
		if version == "" {
			http.Error(w, "version is empty", http.StatusBadRequest)
			return
		}
		versionNum, err := strconv.Atoi(version)
		if err != nil {
			http.Error(w, "version is not number", http.StatusBadRequest)
			return
		}
		if err := m.MigrateTo(r.Context(), uint(versionNum)); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
	})
}

func MigrateUp(mux chi.Router, m migrator) {
	mux.Post("/migrate/up", func(w http.ResponseWriter, r *http.Request) {
		if err := m.MigrateUp(r.Context()); err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
	})
}
