package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"

	"math"
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
	conn, err := net.Dial("tcp", "localhost:9990")
	if err != nil {
		log.Fatalf("dial error: %v", err)
		return
	}
	fmt.Printf("Connected: %#v\n", conn)

	client := &http.Client{
		Transport: &http2.Transport{
			AllowHTTP: true,
			DialTLS: func(network, addr string, cfg *tls.Config) (net.Conn, error) {
				return conn, nil
			},
		},
	}

	go func() {
		req, err := http.NewRequest("POST", "http://localhost:9990/connect", nil)
		if err != nil {
			fmt.Println("error creating request: ", err)
			os.Exit(1)
		}
		resp, err := client.Do(req)
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

		fmt.Println("Response Status: ", string(body))
	}()

	// Provide a http2 server listening to the previously established connection
	h2 := &H2Conn{
		conn: conn,
		server: &http2.Server{
			MaxConcurrentStreams: math.MaxUint32,
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
}
