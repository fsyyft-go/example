package echo

import (
	"bytes"
	"context"
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
		g.Shutdown(context.TODO()) //nolint:errcheck
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
