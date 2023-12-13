package main

import (
	"encoding/json"
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

	r.Get("/", healthCheck)

	fmt.Printf("Listening on port %d\n", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	resp := map[string]string{"result": "ok"}

	byt, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("{ \"error\": \"failed\" }"))
		return
	}

	w.Write(byt)
}
