package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	port = 9999
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Proxy Server!!"))
	})

	fmt.Printf("Listening on port %d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}
