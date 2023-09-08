package main

import (
	exBoot "github.com/fsyyft-go/example/internal/pkg/boot"
)

func main() {
	if err := exBoot.Execute(); nil != err {
		panic(err)
	}
}
