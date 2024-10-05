package main

import (
	"net/http"

	"github.com/IxDay/templ-exp/templates"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		component := templates.Base(templates.Hello("John", nil))
		component.Render(r.Context(), w)
	})
	mux.HandleFunc("/lorem", func(w http.ResponseWriter, r *http.Request) {
		component := templates.Lorem()
		if _, ok := r.Header["Hx-Request"]; !ok {
			component = templates.Base(templates.Hello("John", component))
		}
		component.Render(r.Context(), w)
	})

	http.ListenAndServe(":8080", mux)
}
