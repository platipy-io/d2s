package main

import (
	"log"
	"net/http"

	"github.com/IxDay/templ-exp/app"
	"github.com/IxDay/templ-exp/app/lorem"
)

func main() {
	mux := http.NewServeMux()
	logger := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			log.Println(r.Method + " " + r.URL.Path)
			next(w, r)
		}
	}

	mux.HandleFunc("/", logger(app.Index))
	mux.HandleFunc("/lorem", logger(lorem.Index))

	http.ListenAndServe(":8080", mux)
}
