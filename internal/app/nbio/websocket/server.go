package websocket

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"

	"github.com/lesismal/nbio/nbhttp"

	exSys "github.com/fsyyft-go/example/pkg/sys"
)

/**
 * 参考：https://github.com/lesismal/nbio-examples/blob/master/websocket/server/server.go
 * 		https://github.com/lesismal/nbio/issues/75
 * 关注点：
 * 	1. WebSocket 的 TCP 拆包。
 */

var (
	_ exSys.RunnableWithContext = (*server)(nil)
)

var (
	addr = "127.0.0.1:4444"
)

type (
	server struct {
	}
)

func (s *server) RunContext(ctx context.Context) {
	onDataFrame := false

	mux := &http.ServeMux{}
	mux.HandleFunc("/ws", echoServerHandler(onDataFrame))

	conf := nbhttp.Config{
		Network: "tcp",
		Addrs:   []string{addr},
		Handler: mux,
	}
	svr := nbhttp.NewServer(conf)

	if err := svr.Start(); nil != err {
		slog.Error("WebSocket Start",
			"err", err,
		)
	} else {
		interrupt := make(chan os.Signal, 2)
		signal.Notify(interrupt, os.Interrupt)

		isShutdown := false

		shutdown := func() {
			if !isShutdown {
				isShutdown = true
				if errShutdown := svr.Shutdown(ctx); nil != errShutdown {
					slog.Error("WebSocket Shutdown",
						"err", errShutdown,
					)
				}
			}
		}

		select {
		case <-interrupt:
			shutdown()
		case <-ctx.Done():
			shutdown()
		}
	}
}

func newServer() *server { //nolint:unused
	s := &server{}

	return s
}
