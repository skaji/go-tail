package main

import (
	"context"
	"fmt"

	"github.com/skaji/go-tail"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ch := tail.Run(ctx, tail.NewConfig("access.log"))
	for line := range ch {
		fmt.Printf("got [%s]\n", line)
	}
}
