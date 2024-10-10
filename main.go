package main

import (
	"net/http"

	"github.com/IxDay/templ-exp/app"
	"github.com/IxDay/templ-exp/app/lorem"

	http_ "github.com/IxDay/templ-exp/internal/http"
	"github.com/IxDay/templ-exp/internal/logger"
)

func main() {
	mux := http.NewServeMux()
	log := logger.New(logger.TraceLevel)
	wrapper := func(handler http.HandlerFunc) http.Handler {
		return http_.MiddlewareLogger(log)(http_.MiddlewareRecover(http.HandlerFunc(handler)))
	}
	mux.Handle("/", wrapper(app.Index))
	mux.Handle("/lorem", wrapper(lorem.Index))
	mux.Handle("/panic", wrapper(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("I'm about to panic!")) // this will send a response 200 as we write to resp
		panic("some unknown reason")
	}))
	http.ListenAndServe(":8080", mux)
}
