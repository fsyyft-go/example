package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/lesismal/nbio"
	"github.com/lesismal/nbio/logging"
	"github.com/spf13/cobra"
)

var (
	echoServerCmd *cobra.Command
	echoClientCmd *cobra.Command
)

func init() {
	echoServerCmd = &cobra.Command{
		Use: "echo-server",
		Run: func(cmd *cobra.Command, args []string) {
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
		},
	}

	echoClientCmd = &cobra.Command{
		Use: "echo-client",
		Run: func(cmd *cobra.Command, args []string) {
			var (
				ret    []byte
				buf    = make([]byte, 1024*1024*4)
				addr   = "localhost:8888"
				ctx, _ = context.WithTimeout(context.Background(), 60*time.Second) //nolint:govet
			)

			logging.SetLevel(logging.LevelInfo)

			rand.Read(buf) //nolint:errcheck

			engine := nbio.NewEngine(nbio.Config{})

			done := make(chan int)
			engine.OnData(func(c *nbio.Conn, data []byte) {
				ret = append(ret, data...)
				if len(ret) == len(buf) {
					if bytes.Equal(buf, ret) {
						close(done)
					}
				}
			})

			err := engine.Start()
			if err != nil {
				fmt.Printf("Start failed: %v\n", err)
			}
			defer engine.Stop()

			// net.Dial also can be used here
			c, err := nbio.Dial("tcp", addr)
			if err != nil {
				panic(err)
			}
			engine.AddConn(c) //nolint:errcheck
			c.Write(buf)      //nolint:errcheck

			select {
			case <-ctx.Done():
				logging.Error("timeout")
			case <-done:
				logging.Info("success")
			}
		},
	}

	rootCmd.AddCommand(echoServerCmd, echoClientCmd)
}
