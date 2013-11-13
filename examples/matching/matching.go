package main

import (
	"fmt"
	"os"
	"flag"
	"github.com/progrium/go-plugins"
)

type OutputRenderer struct {
	Match func(string) bool
	Output func() string
}
var OutputRendererExt struct {
	Plugin func(string) OutputRenderer
	Plugins func() []OutputRenderer
}

func main() {
	plugins.LoadFromPath()
	plugins.Register(&OutputRendererExt)

	flag.Parse()

	if flag.Arg(0) == "" {
		fmt.Println("usage: matching <pattern>\n")
		os.Exit(1)
	}

	for _, renderer := range OutputRendererExt.Plugins() {
		if renderer.Match(flag.Arg(0)) {
			fmt.Println(renderer.Output())
			return
		}
	}

}
