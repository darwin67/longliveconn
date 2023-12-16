package main

import (
	"context"
	"net"
)

type Conn struct {
	conn net.Conn
}

type Dialer struct {
	NetDial           func(network, addr string) (net.Conn, error)
	NetDialContext    func(ctx context.Context, network, addr string) (net.Conn, error)
	NetDialTLSContext func(ctx context.Context, network, addr string) (net.Conn, error)
}
