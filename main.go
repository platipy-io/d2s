package main

import (
	"net/http"

	"github.com/IxDay/templ-exp/templates"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		component := templates.Base(templates.Hello("John"))
		component.Render(r.Context(), w)
	})

	http.ListenAndServe(":8080", nil)
}
