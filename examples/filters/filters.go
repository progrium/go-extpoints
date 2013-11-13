package main

import (
	"fmt"
	"os"
	"bufio"
	"github.com/progrium/go-plugins"
)

type TextFilter struct {
	Filter func(string) string
}
var TextFilterExt struct {
	Plugin func(string) TextFilter
	Plugins func() []TextFilter
}

func main() {
	plugins.LoadFromPath()
	
	plugins.Register(&TextFilterExt)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		for _, filter := range TextFilterExt.Plugins() {
			line = filter.Filter(line)
		}
		fmt.Println(line)
	}

}
