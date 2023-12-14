package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"log"
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
	r.Post("/connect", longlive)

	fmt.Printf("Listening on port %d\n", port)

	h2s := &http2.Server{}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: h2c.NewHandler(r, h2s),
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("error running h2 server: %v", err)
	}
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

func longlive(w http.ResponseWriter, r *http.Request) {
	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "server doesn't support hijacking", http.StatusInternalServerError)
		return
	}

	conn, bufrw, err := hj.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	bufrw.WriteString("Now we're speaking raw TCP. Say hi: ")
	bufrw.Flush()
	// s, err := bufrw.ReadString('\n')
	// if err != nil {
	// 	log.Printf("error reading string: %v", err)
	// 	return
	// }

	// fmt.Fprintf(bufrw, "You said: %q\nBye.\n", s)
	// bufrw.Flush()
}
