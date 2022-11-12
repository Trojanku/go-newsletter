package handlers

import (
	"Goo/views"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func FrontPage(mux chi.Router) {
	mux.Get("/", func(w http.ResponseWriter, request *http.Request) {
		tmpl, err := views.LoadTemplate("./views/index.html")
		if err != nil {
			return
		}
		_ = tmpl.Execute(w, nil)
	})
}
