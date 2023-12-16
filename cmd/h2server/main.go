package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"

	// "log"
	"net"
	"net/http"

	"golang.org/x/net/http2"
	// "github.com/go-chi/chi/v5"
	// "github.com/go-chi/chi/v5/middleware"
	// "golang.org/x/net/http2"
	// "golang.org/x/net/http2/h2c"
)

const (
	port = 9999
)

func main() {
	addr := fmt.Sprintf(":%d", port)

	// TCP listener
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

		go echo(conn)
	}

	// APIs
	// r := chi.NewRouter()
	// r.Use(middleware.Logger)

	// r.Get("/", healthCheck)
	// r.Post("/connect", longlive)

	// fmt.Printf("HTTP server listening on port %d\n", port)

	// h2s := &http2.Server{}
	// srv := &http.Server{
	// 	Addr:    fmt.Sprintf(":%d", port),
	// 	Handler: h2c.NewHandler(r, h2s),
	// }

	// if err := srv.Serve(listener); err != nil {
	// 	log.Fatalf("error running h2 server: %v", err)
	// }
}

func echo(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("conn: %#v\n", conn)

	tr := &http2.Transport{
		AllowHTTP: true,
		DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
			return conn, nil
		},
	}
	h2conn, err := tr.NewClientConn(conn)
	if err != nil {
		fmt.Printf("error generating new client conn: %v\n", err)
		return
	}
	fmt.Printf("h2 conn: %#v\n", h2conn.State())

	req, err := http.NewRequest("GET", "http://doesntmatter:3000", nil)
	if err != nil {
		fmt.Printf("error creating request: %#v\n", req)
		return
	}
	resp, err := h2conn.RoundTrip(req)
	if err != nil {
		fmt.Printf("error making request: %#v\n", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error reading body: %#v\n", err)
		return
	}
	fmt.Printf("Response: %s\n", body)

	// fmt.Println("Accepted conn from: ", conn.RemoteAddr())

	// // We can read and write to the connection using conn.Read and conn.Write
	// // For example, echoing the received data back to the client:
	// buffer := make([]byte, 1024)
	// for {
	// 	n, err := conn.Read(buffer)
	// 	if err != nil {
	// 		fmt.Println("error reading from conn: ", err)
	// 		return
	// 	}

	// 	data := buffer[:n]
	// 	fmt.Printf("received data: %s\n", data)

	// 	_, err = conn.Write(data)
	// 	if err != nil {
	// 		fmt.Println("error writing to conn: ", err)
	// 		return
	// 	}
	// }
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
	fmt.Println("Supports hijacking")

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
