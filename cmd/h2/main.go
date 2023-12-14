package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	// "math"
	"net"
	"net/http"
	// "time"

	"golang.org/x/net/http2"
)

func main() {
	fmt.Println("HTTP/2 Client!!")
	// ctx := context.Background()

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
	// if err := h2conn.Ping(ctx); err != nil {
	// 	fmt.Println("ping error: ", err)
	// 	os.Exit(1)
	// }

	// go func() {
	req, err := http.NewRequest("GET", "http://localhost:9999", nil)
	if err != nil {
		fmt.Println("error creating request: ", err)
		os.Exit(1)
	}
	// resp, err := client.Do(req)
	resp, err := h2conn.RoundTrip(req)
	if err != nil {
		fmt.Println("error making request: ", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading body: ", err)
		os.Exit(1)
	}

	fmt.Println("Response: ", string(body))
	// }()

	// // Provide a http2 server listening to the previously established connection
	// h2 := &H2Conn{
	// 	conn: conn,
	// 	server: &http2.Server{
	// 		MaxConcurrentStreams: math.MaxUint32,
	// 	},
	// }

	// if err := h2.Serve(ctx); err != nil {
	// 	log.Fatalf("h2 error: %v", err)
	// }
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
}
