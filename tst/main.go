package main

import (
	"context"
	"fmt"
	"os"
	"tst/fork"
	"tst/orig"
)

func main() {
	fmt.Println("ORIG:")
	if err := orig.WithFmt(42).Render(context.Background(), os.Stdout); err != nil {
		panic(err)
	}
	fmt.Println("")

	fmt.Println("FORK:")
	if err := fork.WithFmt(42).Render(context.Background(), os.Stdout); err != nil {
		panic(err)
	}
	fmt.Println("")
}
