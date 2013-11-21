package main

import (
	"fmt"
	"github.com/progrium/go-plugins"
	"github.com/progrium/go-plugins/ottojs"
)

type ProgramObserver struct {
	ProgramStarted func()
	ProgramEnded   func()
}

var ProgramObserverExt struct {
	Plugin  func(string) ProgramObserver
	Plugins func() []ProgramObserver
}

type MyStaticPlugin struct{}

func (p MyStaticPlugin) ProgramStarted() {
	fmt.Println(plugins.GetGlobal("prefix").(string) + ": start")
}

func (p MyStaticPlugin) ProgramEnded() {
	fmt.Println(plugins.GetGlobal("prefix").(string) + ": end")
}

func main() {
	plugins.RegisterRuntime(ottojs.GetRuntime())
	plugins.ExtensionPoint(&ProgramObserverExt)
	plugins.SetGlobals(map[string]interface{}{
		"prefix": "Static plugin",
	})

	plugins.StaticPlugin(&MyStaticPlugin{}, []string{
		"ProgramObserver",
	})

	for _, observer := range ProgramObserverExt.Plugins() {
		observer.ProgramStarted()
	}

	fmt.Println("Hello World")

	for _, observer := range ProgramObserverExt.Plugins() {
		observer.ProgramEnded()
	}
}
