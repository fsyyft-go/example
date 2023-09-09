package netstd

import (
	"bytes"
	"context"
	"crypto/rand"
	"log/slog"
	"net"
	"time"

	"github.com/lesismal/nbio"

	exSys "github.com/fsyyft-go/example/pkg/sys"
)

var (
	_ exSys.Runnable = (*client)(nil)
)

type (
	client struct {
	}
)

func (c *client) Run() {
	var (
		ret    []byte
		buf    = make([]byte, 32)
		ctx, _ = context.WithTimeout(context.Background(), time.Second) //nolint:govet
	)

	rand.Read(buf) //nolint:errcheck

	g := nbio.NewGopher(nbio.Config{})

	done := make(chan int)
	g.OnData(func(c *nbio.Conn, data []byte) {
		ret = append(ret, data...)
		if len(ret) == len(buf) {
			if bytes.Equal(buf, ret) {
				close(done)
			}
		}
	})

	err := g.Start()
	if err != nil {
		slog.Error("客户端 启动连接失败", "err", err)
	}
	defer g.Stop()

	cli, err := net.Dial("tcp", addr)
	if err != nil {
		slog.Error("客户端 建议连接失败", "err", err)
	}
	g.AddConn(cli) //nolint:errcheck
	go func() {
		cli.Write(buf) //nolint:errcheck
		slog.Info("客户端 数据发送完成。")
		done <- 0
	}()

	select {
	case <-ctx.Done():
		slog.Error("客户端 执行超时。")
	case <-done:
		slog.Info("客户端 执行完成。")
	}
}

func NewClient() *client {
	c := &client{}

	return c
}
