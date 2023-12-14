package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/http2"
)

func main() {
	fmt.Println("HTTP/2 Client!!")

	ctx := context.Background()
	dialCtx, dialCancel := context.WithTimeout(ctx, 10*time.Second)
	defer dialCancel()

	// Establish TCP connection with remote server
	dialer := net.Dialer{}
	conn, err := dialer.DialContext(dialCtx, "tcp", "localhost:9999")
	if err != nil {
		log.Fatalf("dial error: %v", err)
		return
	}
	fmt.Printf("Connected: %#v\n", conn)

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
