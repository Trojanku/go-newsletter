package handlers

import (
	"Goo/views"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func FrontPage(mux chi.Router) {
	mux.Get("/", func(w http.ResponseWriter, request *http.Request) {
		tmpl, err := views.LoadTemplate()
		if err != nil {
			fmt.Printf("Error loading template: %v \n", err)
			return
		}
		_ = tmpl.Execute(w, nil)
	})
}
