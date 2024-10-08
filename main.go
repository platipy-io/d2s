package main

import (
	"log"
	"net/http"

	"github.com/IxDay/templ-exp/templates"
)

func main() {
	mux := http.NewServeMux()
	logger := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.Method + " " + r.URL.Path)
			next(w, r)
		}
	}

	mux.HandleFunc("/", logger(func(w http.ResponseWriter, r *http.Request) {
		component := templates.Base(templates.Hello("John", nil))
		component.Render(r.Context(), w)
	}))
	mux.HandleFunc("/lorem", logger(func(w http.ResponseWriter, r *http.Request) {
		component := templates.Lorem()
		if _, ok := r.Header["Hx-Request"]; !ok {
			component = templates.Base(templates.Hello("John", component))
		}
		component.Render(r.Context(), w)
	}))

	http.ListenAndServe(":8080", mux)
}
