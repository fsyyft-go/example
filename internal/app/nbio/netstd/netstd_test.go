package netstd

import (
	"testing"
	"time"
)

func TestNetstd(t *testing.T) {
	go func() {
		s := NewServer()
		s.Run()
	}()
	time.Sleep(10 * time.Millisecond)
	c := NewClient()
	c.Run()
}

func TestMain(m *testing.M) {
	m.Run()
}
