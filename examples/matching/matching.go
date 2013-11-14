package main

import (
	"fmt"
	"os"
	"flag"
	"github.com/progrium/go-plugins"
	"github.com/progrium/go-plugins/ottojs"
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
	plugins.RegisterRuntime(ottojs.GetRuntime())
	plugins.LoadFromPath()
	plugins.ExtensionPoint(&OutputRendererExt)

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
