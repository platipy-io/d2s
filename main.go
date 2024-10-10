package main

import (
	"net/http"

	"github.com/IxDay/templ-exp/app"
	"github.com/IxDay/templ-exp/app/lorem"
	"github.com/IxDay/templ-exp/internal/logger"
)

func main() {
	mux := http.NewServeMux()
	log := logger.New(logger.TraceLevel)
	mux.Handle("/", logger.Middleware(log)(http.HandlerFunc(app.Index)))
	mux.Handle("/lorem", logger.Middleware(log)(http.HandlerFunc(lorem.Index)))
	http.ListenAndServe(":8080", mux)
}
