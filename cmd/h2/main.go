package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	// "time"

	"golang.org/x/net/http2"
)

func main() {
	fmt.Println("HTTP/2 Client!!")
	ctx := context.Background()

	// Establish TCP connection with remote server
	conn, err := net.Dial("tcp", "localhost:9999")
	if err != nil {
		log.Fatalf("dial error: %v", err)
		return
	}
	fmt.Printf("Connected: %#v\n", conn)

	tr := &http2.Transport{
		AllowHTTP: true,
		DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
			return conn, nil
		},
	}
	h2conn, err := tr.NewClientConn(conn)
	if err != nil {
		fmt.Println("error creating new h2 conn: ", err)
		os.Exit(1)
	}

	fmt.Printf("h2 conn: %#v\n", h2conn.State())

	// go func() {
	// 	req, err := http.NewRequest("GET", "http://localhost:9999", nil)
	// 	if err != nil {
	// 		fmt.Println("error creating request: ", err)
	// 		os.Exit(1)
	// 	}
	// 	// resp, err := client.Do(req)
	// 	resp, err := h2conn.RoundTrip(req)
	// 	if err != nil {
	// 		fmt.Println("error making request: ", err)
	// 		os.Exit(1)
	// 	}
	// 	defer resp.Body.Close()

	// 	body, err := ioutil.ReadAll(resp.Body)
	// 	if err != nil {
	// 		fmt.Println("error reading body: ", err)
	// 		os.Exit(1)
	// 	}

	// 	fmt.Println("Response: ", string(body))

	// 	fmt.Printf("Conn: %#v\n", h2conn.State())
	// }()

	// Provide a http2 server listening to the previously established connection
	h2 := &H2Conn{
		conn: conn,
		server: &http2.Server{
			MaxConcurrentStreams: 2,
		},
	}

	if err := h2.Serve(ctx); err != nil {
		log.Fatalf("h2 error: %v", err)
	}
}

type H2Conn struct {
	conn   net.Conn
	server *http2.Server
}

func (c *H2Conn) Serve(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		c.conn.Close()
	}()

	c.server.ServeConn(c.conn, &http2.ServeConnOpts{
		Context: ctx,
		Handler: c,
	})

	return nil
}

func (c *H2Conn) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// w.Header().Add("Content-Type", "application/json")

	fmt.Printf("Req: %#v\n", r)

	resp, err := json.Marshal(map[string]string{"result": "ok"})
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("error serializing to json"))
		return
	}

	w.WriteHeader(200)
	w.Write(resp)
}
