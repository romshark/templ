package main

import (
	"context"
	"fmt"
	"os"
	"tst/fork"
	"tst/orig"

	"github.com/a-h/templ"
)

func main() {
	const id = 42

	printComp("orig", orig.WithFmt(id))
	printComp("orig-concat", orig.Concat(id))
	printComp("orig-strbuild", orig.Strbuild(id))
	printComp("fork", fork.WithFmt(id))
}

func printComp(name string, comp templ.Component) {
	fmt.Println(name + ":")
	if err := comp.Render(context.Background(), os.Stdout); err != nil {
		panic(err)
	}
	fmt.Println("")
}
