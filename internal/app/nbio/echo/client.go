package echo

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/lesismal/nbio"
	"github.com/lesismal/nbio/logging"

	exSys "github.com/fsyyft-go/example/pkg/sys"
)

var (
	_ exSys.Runnable = (*client)(nil)
)

type (
	client struct { //nolint:unused
	}
)

func (c *client) Run() { //nolint:unused
	var (
		ret    []byte
		buf    = make([]byte, 1024*1024*4)
		addr   = "127.0.0.1:44444"
		ctx, _ = context.WithTimeout(context.Background(), 60*time.Second) //nolint:govet
	)

	logging.SetLevel(logging.LevelInfo)

	rand.Read(buf) //nolint:errcheck

	engine := nbio.NewEngine(nbio.Config{})

	done := make(chan int)
	engine.OnData(func(c *nbio.Conn, data []byte) {
		ret = append(ret, data...)
		if len(ret) == len(buf) {
			if bytes.Equal(buf, ret) {
				close(done)
			}
		}
	})

	err := engine.Start()
	if err != nil {
		fmt.Printf("Start failed: %v\n", err)
	}
	defer engine.Stop()

	// net.Dial also can be used here
	cl, err := nbio.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	engine.AddConn(cl) //nolint:errcheck,typecheck
	cl.Write(buf)      //nolint:errcheck

	select {
	case <-ctx.Done():
		logging.Error("timeout")
	case <-done:
		logging.Info("success")
	}
}

func NewClient() *client {
	c := &client{}

	return c
}
