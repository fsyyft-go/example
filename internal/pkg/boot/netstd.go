package boot

import (
	"github.com/spf13/cobra"

	"github.com/fsyyft-go/example/internal/app/nbio/netstd"
)

var (
	netstdServerCmd *cobra.Command
	netstdClientCmd *cobra.Command
)

func init() {
	netstdServerCmd = &cobra.Command{
		Use: "netstd-server",
		Run: func(cmd *cobra.Command, args []string) {
			s := netstd.NewServer()
			s.Run()
		},
	}

	netstdClientCmd = &cobra.Command{
		Use: "netstd-client",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	rootCmd.AddCommand(netstdServerCmd, netstdClientCmd)
}
