package main

import (
	"fmt"
	"os"
	"bufio"
	"github.com/progrium/go-plugins"
	"github.com/progrium/go-plugins/ottojs"
)

type TextFilter struct {
	Filter func(string) string
}
var TextFilterExt struct {
	Plugin func(string) TextFilter
	Plugins func() []TextFilter
}

func main() {
	plugins.RegisterRuntime(ottojs.GetRuntime())
	plugins.LoadFromPath()

	plugins.ExtensionPoint(&TextFilterExt)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		for _, filter := range TextFilterExt.Plugins() {
			line = filter.Filter(line)
		}
		fmt.Println(line)
	}

}
