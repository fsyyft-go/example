package echo

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"testing"
	"time"

	"github.com/lesismal/nbio"

	exTesting "github.com/fsyyft-go/example/pkg/testing"
)

const (
	addr = "127.0.0.1:4444"
)

var (
	buf = make([]byte, 6*1024*1024)
)

func init() {
	rand.Read(buf) //nolint:errcheck
}

func writeComplete(c *nbio.Conn, data []byte) (int, error) {
	offset := 0
	msgLen := len(data)
	for {
		n, err := c.Write(data[offset:])
		fmt.Printf("write %d %s\n", n, err)
		offset += n
		if err != nil || offset == msgLen {
			return offset, err
		}
		time.Sleep(time.Millisecond * 500)
	}

}

func testServer(ready chan error) error {
	g := nbio.NewGopher(nbio.Config{
		Network:            "tcp",
		Addrs:              []string{addr},
		MaxWriteBufferSize: 6 * 1024 * 1024,
	})

	g.OnOpen(func(c *nbio.Conn) {
		_, err := writeComplete(c, buf)
		if err != nil {
			fmt.Printf("write failed: %s\n", err)
		}
	})
	/**
	 * 2023/09/14 12:15:41.483 [INF] NBIO[NB] stop
	 * 2023/09/14 12:15:41.483 [INF] NBIO[NB] stop
	 * 2023/09/14 12:15:41.483 [ERR] Timer[NB] exec call failed: sync: negative WaitGroup counter
	 * goroutine 29 [running]:
	 * github.com/lesismal/nbio/timer.(*Timer).Async.func1.1()
	 * 	/home/runner/go/pkg/mod/github.com/lesismal/nbio@v1.3.18/timer/timer.go:69 +0x72
	 * panic({0x5538c0?, 0x5c0e40?})
	 * 	/opt/hostedtoolcache/go/1.21.1/x64/src/runtime/panic.go:914 +0x21f
	 * sync.(*WaitGroup).Add(0x0?, 0x0?)
	 * 	/opt/hostedtoolcache/go/1.21.1/x64/src/sync/waitgroup.go:62 +0xd8
	 * sync.(*WaitGroup).Done(0xc0000da000?)
	 * 	/opt/hostedtoolcache/go/1.21.1/x64/src/sync/waitgroup.go:87 +0x1a
	 * github.com/fsyyft-go/example/internal/app/nbio/echo.testServer.(*Engine).OnClose.func5.1()
	 * 	/home/runner/go/pkg/mod/github.com/lesismal/nbio@v1.3.18/engine.go:239 +0x71
	 * github.com/lesismal/nbio/timer.(*Timer).Async.func1()
	 * 	/home/runner/go/pkg/mod/github.com/lesismal/nbio@v1.3.18/timer/timer.go:73 +0x4d
	 * created by github.com/lesismal/nbio/timer.(*Timer).Async in goroutine 28
	 * 	/home/runner/go/pkg/mod/github.com/lesismal/nbio@v1.3.18/timer/timer.go:63 +0x67
	 */
	g.OnClose(func(c *nbio.Conn, err error) {
		if r := recover(); nil != r {
			exTesting.Printf("OnClose 发生异常：%[1]s", r)
		}
		g.Stop()
	})

	err := g.Start()
	if err != nil {
		return fmt.Errorf("nbio.Start failed: %w", err)
	}
	ready <- err

	go func() {
		if r := recover(); nil != r {
			exTesting.Printf("关闭 groutine 发生异常：%[1]s", r)
		}

		time.Sleep(time.Second * 3)
		g.Stop()
	}()
	// defer g.Stop()

	g.Wait()
	return nil
}

func testClient(msgLen int) error {
	var (
		ret  []byte
		addr = addr
	)
	c, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println(err)
		return err
	}

	i := 0
	line := make([]byte, 60000)
	for {

		n, err := c.Read(line)
		if err != nil && !errors.Is(err, io.EOF) {
			return fmt.Errorf("error read: %d %w", n, err)
		}
		if errors.Is(err, io.EOF) {
			time.Sleep(time.Second * 5)
		}
		i++
		ret = append(ret, line[:n]...)
		fmt.Printf("client received %d %d %d of %d\n", i, n, len(ret), len(buf))
		if len(ret) == len(buf) {
			if bytes.Equal(buf, ret) {
				return nil
			}
			return fmt.Errorf("ret, does not match buf")
		}

	}
}
func Test_main(t *testing.T) {
	ready := make(chan error)
	go func() {
		err := testServer(ready)
		if err != nil {
			log.Fatal(err)
		}
	}()

	err := <-ready
	if err != nil {
		t.Fatal(err)
	}

	err = testClient(1024 * 1024 * 4)
	if err != nil {
		t.Fatal(err)
	}
}
