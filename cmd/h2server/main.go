package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

const (
	tcpPort  = 9999
	httpPort = 9990
)

func main() {
	addr := fmt.Sprintf(":%d", tcpPort)

	// TCP listener
	go func() {
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			fmt.Println("error listening to TCP: ", err)
			return
		}
		defer listener.Close()
		fmt.Println("TCP listener on: ", listener.Addr())

		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("error accepting conn: ", err)
				continue
			}

			go handleConnection(conn)
		}
	}()

	// APIs
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", healthCheck)
	r.Post("/connect", longlive)

	fmt.Printf("HTTP server listening on port %d\n", httpPort)

	h2s := &http2.Server{}
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", httpPort),
		Handler: h2c.NewHandler(r, h2s),
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("error running h2 server: %v", err)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Accepted conn from: ", conn.RemoteAddr())

	// We can read and write to the connection using conn.Read and conn.Write
	// For example, echoing the received data back to the client:
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("error reading from conn: ", err)
			return
		}

		data := buffer[:n]
		fmt.Printf("received data: %s\n", data)

		_, err = conn.Write(data)
		if err != nil {
			fmt.Println("error writing to conn: ", err)
			return
		}
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
