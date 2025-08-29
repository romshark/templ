package main

import (
	"context"
	"fmt"
	"os"
)

func main() {
	if err := WithFmt(42).Render(context.Background(), os.Stdout); err != nil {
		panic(err)
	}
	fmt.Println("")
}
