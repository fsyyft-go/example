package testing

import (
	"fmt"
)

var (
	// logHeader 日志输出前缀。
	//
	// 使用 var 而不是使用 const，可以根据需要在构建时重新赋值。
	logHeader = "-=>       "
)

func Println(a ...interface{}) {
	fmt.Print(logHeader)
	// fmt.Print(a...) 所有参数连在一起。
	// fmt.Println(a...) 参数之间空格分割。
	fmt.Println(a...)
}

func Printf(format string, a ...interface{}) {
	fmt.Print(logHeader)
	fmt.Printf(format, a...)
}
