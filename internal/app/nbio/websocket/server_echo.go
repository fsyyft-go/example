package websocket

import (
	"log/slog"
	"net/http"

	"github.com/lesismal/nbio/nbhttp/websocket"
)

type (
	echoServer struct { //nolint:unused
		upgrader    *websocket.Upgrader
		onDataFrame bool
	}
)

func (s *echoServer) message(c *websocket.Conn, messageType websocket.MessageType, data []byte) { //nolint:unused
	// 数据原样回写。
	if err := c.WriteMessage(messageType, data); nil != err {
		slog.Error("WebSocket OnMessage",
			"err", err,
		)
	}
}

func (s *echoServer) open(c *websocket.Conn) { //nolint:unused
	slog.Info("WebSocket OnOpen",
		"RemoteAddr", c.RemoteAddr().String(),
	)
}

func (s *echoServer) close(c *websocket.Conn, err error) { //nolint:unused
	slog.Info("WebSocket OnClose",
		"RemoteAddr", c.RemoteAddr().String(),
		"err", err,
	)
}

func echoServerUpgrader(onDataFrame bool) *websocket.Upgrader { //nolint:unused
	u := websocket.NewUpgrader()
	es := &echoServer{
		upgrader:    u,
		onDataFrame: onDataFrame,
	}

	u.OnOpen(es.open)
	u.OnClose(es.close)
	u.OnMessage(es.message)

	return u
}

func echoServerHandler(onDataFrame bool) func(http.ResponseWriter, *http.Request) { //nolint:unused
	u := echoServerUpgrader(false)

	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = u.Upgrade(w, r, nil)
	}
}
