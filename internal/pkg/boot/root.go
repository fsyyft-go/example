package boot

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "nbio",
	Short: "这是一个关于 nbio 库的示例。",
	Long:  `这是一个关于 nbio 库的示例，也可以通过官方示例代码中查看：https://github.com/lesismal/nbio-examples。`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(args)
	},
}

func init() {
}

func Execute() error {
	return rootCmd.Execute()
}
