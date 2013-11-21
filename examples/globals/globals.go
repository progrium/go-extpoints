package main

import (
	"fmt"
	"github.com/progrium/go-plugins"
	"github.com/progrium/go-plugins/ottojs"
)

type GlobalsUser struct {
	UseGlobals func()
}

var GlobalsUserExt struct {
	Plugin  func(string) GlobalsUser
	Plugins func() []GlobalsUser
}

func main() {
	plugins.RegisterRuntime(ottojs.GetRuntime())
	plugins.LoadFromPath()
	plugins.SetGlobals(map[string]interface{}{
		"Foo.Bar.Baz": "Hello world",
		"Foo.Log": func(text string, n int) {
			fmt.Println(text)
			fmt.Println(n)
		},
		"Foo.Bar.Qux": 42,
	})

	plugins.ExtensionPoint(&GlobalsUserExt)

	for _, obj := range GlobalsUserExt.Plugins() {
		obj.UseGlobals()
	}

}
