package echo

import (
	"github.com/fsyyft-go/example/pkg/sys"
)

var (
	_ sys.Runnable = (*client)(nil)
)

type (
	client struct { //nolint:unused
	}
)

func (c *client) Run() { //nolint:unused
}

func NewClient() *client {
	c := &client{}

	return c
}
