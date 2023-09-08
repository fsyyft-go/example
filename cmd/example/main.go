package main

import (
	"fmt"
)

func main() {
	if err := run(); nil != err {
		panic(err)
	}
}

func run() error {
	fmt.Println("这是一个示例")
	return nil
}
