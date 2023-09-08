package netstd

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/lesismal/nbio"

	exSys "github.com/fsyyft-go/example/pkg/sys"
)

var (
	_ exSys.Runnable = (*server)(nil)
)

var (
	addr = "127.0.0.1:4444"
)

type (
	server struct { //nolint:unused
	}
)

func (s *server) Run() { //nolint:unused
	g := nbio.NewGopher(nbio.Config{})
	g.OnOpen(func(c *nbio.Conn) {
		slog.Info("OnOpen:", c.RemoteAddr().String()) //nolint:govet
	})

	g.OnData(func(c *nbio.Conn, data []byte) {
		c.Write(append([]byte{}, data...)) //nolint:errcheck
	})

	err := g.Start()
	if err != nil {
		fmt.Printf("nbio.Start failed: %v\n", err)
		return
	}
	defer g.Stop()

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		slog.Error(err.Error())
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			slog.Info("Accept failed:", err)
			continue
		}
		g.AddConn(conn) //nolint:errcheck
	}
}

func NewServer() *server {
	s := &server{}

	return s
}
