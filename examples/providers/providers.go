package main

import (
	"fmt"
	"github.com/progrium/go-plugins"
	"github.com/progrium/go-plugins/ottojs"
)

type ListProvider struct {
	GetItems func() []interface{}
}
var ListProviderExt struct {
	Plugin func(string) ListProvider
	Plugins func() []ListProvider
}

func main() {
	plugins.RegisterRuntime(ottojs.GetRuntime())
	plugins.LoadFromPath()
	
	plugins.ExtensionPoint(&ListProviderExt)

	fmt.Println("Here is a list, produced by plugins:")

	for _, provider := range ListProviderExt.Plugins() {
		for _, item := range provider.GetItems() {
			fmt.Println(" * " + item.(string))
		}
	}

}
