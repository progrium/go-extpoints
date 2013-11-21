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

func main() {
	plugins.RegisterRuntime(ottojs.GetRuntime())
	plugins.LoadFromPath()

	plugins.ExtensionPoint(&ProgramObserverExt)

	for _, observer := range ProgramObserverExt.Plugins() {
		observer.ProgramStarted()
	}

	fmt.Println("Hello World")

	for _, observer := range ProgramObserverExt.Plugins() {
		observer.ProgramEnded()
	}
}
