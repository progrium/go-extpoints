package main

import (
	"fmt"
	"os"
	"flag"
	"github.com/progrium/go-plugins"
)

type PluggableBackend struct {
	Name func() string
	Process func(string, string) string
}
var PluggableBackendExt struct {
	Plugin func(string) PluggableBackend
	Plugins func() []PluggableBackend
}

func main() {
	plugins.LoadFromPath()
	plugins.Register(&PluggableBackendExt)

	flag.Parse()

	if flag.Arg(0) == "" {
		fmt.Println("usage: backends <backend> <arg1> <arg2>\n")
		fmt.Println("available backends:")
		for _, backend := range PluggableBackendExt.Plugins() {
			fmt.Println("  "+backend.Name())
		}
		fmt.Println("")
		os.Exit(1)
	}

	backend := PluggableBackendExt.Plugin(flag.Arg(0))
	fmt.Println(backend.Process(flag.Arg(1), flag.Arg(2)))

}
