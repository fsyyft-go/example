package boot

import (
	"github.com/spf13/cobra"

	"github.com/fsyyft-go/example/internal/app/nbio/echo"
)

var (
	echoServerCmd *cobra.Command
	echoClientCmd *cobra.Command
)

func init() {
	echoServerCmd = &cobra.Command{
		Use: "echo-server",
		Run: func(cmd *cobra.Command, args []string) {
			s := echo.NewServer()
			s.Run()
		},
	}

	echoClientCmd = &cobra.Command{
		Use: "echo-client",
		Run: func(cmd *cobra.Command, args []string) {
			c := echo.NewClient()
			c.Run()
		},
	}

	rootCmd.AddCommand(echoServerCmd, echoClientCmd)
}
