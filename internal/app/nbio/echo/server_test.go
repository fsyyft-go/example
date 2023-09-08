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

func writeComplete(c *nbio.Conn, data []byte) (int, error) {
	offset := 0
	msgLen := len(data)
	endIndex := 0

	i := 0
	for {
		endIndex = offset + (1024 * 1024)
		if endIndex > msgLen {
			endIndex = msgLen - 1
		}
		n, err := c.Write(data[offset:endIndex])

		i++
		exTesting.Printf("服务端 发送数据：%[1]d %[2]d %[3]d %[4]d %[5]v\n", i, n, offset, endIndex, err)
		offset += n
		if err != nil || offset == msgLen {
			exTesting.Printf("服务端 发送结束：%[1]d %[2]v\n", offset, err)
			return offset, err
		}
		time.Sleep(time.Millisecond * 10)
	}

}

func testServer(ready chan error, exit chan interface{}) error {
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

	g.OnClose(func(c *nbio.Conn, err error) {
		if r := recover(); nil != r {
			exTesting.Printf("服务端 OnClose 发生错误：%[1]v\n", r)
		}
		exTesting.Printf("服务端 OnClose 开始执行：%[1]v\n", err)
		exit <- struct{}{}
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

func testClient(msgLen int, exit chan interface{}) error {
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
		select {
		case <-exit:
			// TODO 这种情况，不是预期的，需要找出来原因。
			exTesting.Printf("客户端 因服务端异常而退出\n")

			return nil
		default:
			n, err := c.Write([]byte{})
			if nil != err {
				return fmt.Errorf("客户端 发送数据（空包）异常：%[1]s\n", err.Error())
			}

			n, err = c.Read(line)

			if nil != err && !errors.Is(err, io.EOF) {
				return fmt.Errorf("error read: %d %d %w", i, n, err)
			} else if nil != err && errors.Is(err, io.EOF) {
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
}

func TestWriteData(t *testing.T) {
	/**
	 * 1. 对需要进行发送的数据，进行随机的填充。
	 * 2. 异步启动服务端。
	 * 3. 同步等服务端启动完成后，启动客户端。
	 * 4. 客户端连接后，服务端一次性发送消息。
	 * 5. 客户端收到所有消息后，客户端断开连接。
	 * 6. 客户端断开连接时，服务端的回调函数触发，关闭服务。
	 * 7. 服务端退出后，测试结束。
	 */
	assertions := assert.New(t)

	// 对 buf 进行随机的填充。
	rand.Read(buf) //nolint:errcheck

	exit := make(chan interface{})

	ready := make(chan error)
	go func() {
		err := testServer(ready, exit)
		if err != nil {
			log.Fatal(err)
		}
	}()

	err := <-ready
	assertions.Nil(err)

	err = testClient(1024*1024*4, exit)
	time.Sleep(500 * time.Millisecond)
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
