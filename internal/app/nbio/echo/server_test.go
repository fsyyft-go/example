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
	"sync"
	"testing"
	"time"

	"github.com/lesismal/nbio"
	"github.com/stretchr/testify/assert"

	exTesting "github.com/fsyyft-go/example/pkg/testing"
)

const (
	addr1 = "127.0.0.1:4444"
	addr2 = "127.0.0.1:4445"
)

var (
	buf = make([]byte, 6*1024*1024)
)

func init() {
	// 对 buf 进行随机的填充。
	rand.Read(buf) //nolint:errcheck
}

func writeComplete(c *nbio.Conn, data []byte) (int, error) {
	offset := 0
	msgLen := len(data)
	for {
		n, err := c.Write(data[offset:])
		exTesting.Printf("服务端 发送数据：%[1]d %[2]v\n", n, err)
		offset += n
		if err != nil || offset == msgLen {
			exTesting.Printf("服务端 发送结束：%[1]d %[2]v\n", offset, err)
			return offset, err
		}
		time.Sleep(time.Millisecond * 100)
	}

}

func testServer(ready chan error) error {
	g := nbio.NewGopher(nbio.Config{
		Network:            "tcp",
		Addrs:              []string{addr1},
		MaxWriteBufferSize: 6 * 1024 * 1024,
	})

	g.OnOpen(func(c *nbio.Conn) {
		_, err := writeComplete(c, buf)
		if err != nil {
			exTesting.Printf("服务端 发送数据发生错误：%[1]s\n", err.Error())
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
			exTesting.Printf("服务端 OnClose 发生错误：%[1]v\n", r)
		}
		g.Stop()
	})

	err := g.Start()
	if err != nil {
		exTesting.Printf("服务端 Start 发生错误：%[1]v\n", err)
	}
	ready <- err

	defer func() {
		if r := recover(); nil != r {
			exTesting.Printf("服务端 defer 发生错误：%[1]v\n", r)
		}

		exTesting.Printf("服务端 defer 执行")

		g.Stop()
	}()

	g.Wait()

	return nil
}

func testClient(msgLen int) error {
	var (
		ret  []byte
		addr = addr1
	)
	c, err := net.Dial("tcp", addr)
	if err != nil {
		exTesting.Printf("客户端 建立连接 失败：%[1]s\n", err.Error())
		return err
	}

	i := 0
	line := make([]byte, 600000)
	for {
		n, err := c.Read(line)
		if 0 == n || n >= 65535 || (err != nil && !errors.Is(err, io.EOF)) {
			return fmt.Errorf("error read: %d %w", n, err)
		}
		if errors.Is(err, io.EOF) {
			time.Sleep(time.Millisecond * 50)
		}
		i++
		ret = append(ret, line[:n]...)
		exTesting.Printf("客户端 收到数据：%d %d %d of %d\n", i, n, len(ret), len(buf))
		if len(ret) == len(buf) {
			if bytes.Equal(buf, ret) {
				return nil
			}
			return fmt.Errorf("ret, does not match buf")
		}

	}
}

func TestWriteData(t *testing.T) {
	assertions := assert.New(t)

	ready := make(chan error)
	go func() {
		err := testServer(ready)
		if err != nil {
			log.Fatal(err)
		}
	}()

	err := <-ready
	assertions.Nil(err)

	err = testClient(1024 * 1024 * 4)
	assertions.Nil(err)
}

type myServer struct {
	nbio.Engine
	sync.Once
}

func (s *myServer) Stop() {
	s.Once.Do(func() {
		s.Engine.Stop()
	})
}

func (s *myServer) Shutdown(ctx context.Context) error {
	var err error

	s.Once.Do(func() {
		err = s.Engine.Shutdown(ctx)
	})

	return err
}

func newMyServer(conf nbio.Config) *myServer {
	g := &myServer{
		Engine: *nbio.NewGopher(conf),
	}

	return g
}

func TestServerStop(t *testing.T) {
	t.Run("Stop", func(t *testing.T) {
		g := newMyServer(nbio.Config{
			Network:            "tcp",
			Addrs:              []string{addr2},
			MaxWriteBufferSize: 6 * 1024 * 1024,
		})

		g.OnOpen(func(c *nbio.Conn) {
			exTesting.Printf("ServerStop 方法的 OnOpen：%[1]s\n", c.RemoteAddr())
		})

		g.OnClose(func(c *nbio.Conn, err error) {
			exTesting.Printf("ServerStop 方法的 OnClose：%[1]s\n", c.RemoteAddr())
		})

		g.OnStop(func() {
			exTesting.Printf("ServerStop 方法的 OnStop\n")
		})

		if err := g.Start(); nil != err {
			exTesting.Printf("ServerStop 方法的 Start 发生错误：%[1]s\n", err.Error())
		} else {
			time.Sleep(time.Millisecond * 100)
			g.Stop()
			_ = g.Shutdown(context.TODO())
			g.Stop()
			_ = g.Shutdown(context.TODO())
			g.Stop()
			_ = g.Shutdown(context.TODO())
		}
	})
}

func Test_main(t *testing.T) {

}
