package echo

import (
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
}

func NewServer() *server {
	s := &server{}

	return s
}
