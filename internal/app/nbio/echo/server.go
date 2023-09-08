package echo

import (
	"fmt"

	"github.com/lesismal/nbio"

	"github.com/fsyyft-go/example/pkg/sys"
)

var (
	_ sys.Runnable = (*server)(nil)
)

type (
	server struct { //nolint:unused
	}
)

func (s *server) Run() { //nolint:unused

	engine := nbio.NewEngine(nbio.Config{
		Network:            "tcp",
		Addrs:              []string{"127.0.0.1:44444"},
		MaxWriteBufferSize: 6 * 1024 * 1024,
	})

	engine.OnData(func(c *nbio.Conn, data []byte) {
		c.Write(append([]byte{}, data...)) //nolint:errcheck
	})

	err := engine.Start()
	if err != nil {
		fmt.Printf("nbio.Start failed: %v\n", err)
		return
	}
	defer engine.Stop()

	<-make(chan int)
}

func NewServer() *server {
	s := &server{}

	return s
}
