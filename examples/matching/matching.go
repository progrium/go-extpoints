package main

import (
	"fmt"
	"os"
	"flag"
	"github.com/progrium/go-plugins"
)

type OutoutRenderer struct {
	Match func(string) bool
	Output func() string
}
var OutoutRendererExt struct {
	Plugin func(string) OutoutRenderer
	Plugins func() []OutoutRenderer
}

func main() {
	plugins.LoadFromPath()
	plugins.Register(&OutoutRendererExt)

	flag.Parse()

	if flag.Arg(0) == "" {
		fmt.Println("usage: matching <pattern>\n")
		os.Exit(1)
	}

	for _, renderer := range OutoutRendererExt.Plugins() {
		if renderer.Match(flag.Arg(0)) {
			fmt.Println(renderer.Output())
			return
		}
	}

}
