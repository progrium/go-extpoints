# go-extpoints

This package, named short for "Go extension points", provides a system for making your Go packages, libraries, and applications extensible in a standard way. 

It expands on the foundation of Go interfaces and provides a meta-API for accessing and registering "extensions", objects that implement "extension point" interfaces. It not only lets third-party packages hook in as build-time extensions, but also encourages better organization of your own package and subpackages. 

## Getting the tool

	$ go install github.com/progrium/go-extpoints

## Quick Example

Here is a simple Go application that turns its `main()` into an extension point based on a `ProgramParticipant` interface.

```
// file main.go
//go:generate go-extpoints
package main

import (
	"./extpoints"
)

func main() {
	for _, extension := range extpoints.ProgramParticipants.All() {
		extension.Main()
	}
}
```

We create an `extpoints` subpackage with a Go file in it to define our extension point interface:

```
// file extpoints/interfaces.go
package extpoints

type ProgramParticipant interface {
	Main()
}
```

This is where `go-extpoints` comes in. It hooks into `go generate` to produce extension point code in our `extpoints` subpackage around the interfaces defined there.

	$ go generate
	 ....
	$ go install

But it doesn't *do* anything! Let's make a builtin `ProgramParticipant`.

```
// file builtin.go
package main

import (
	"fmt"
	"./extpoints"
)

func init() {
	extpoints.Register(new(BuiltinParticipant))
}

type BuiltinParticipant struct {}

func (p *BuiltinParticipant) Main() {
	fmt.Println("Hello world!")
}
```

Now when we build and run the app, it does something! This is a trivial example and the value of extension points becomes clearer with larger applications. But just consider now that `builtin.go` *could* exist in another package in another repo, and you'd just have to import it and rebuild to let it hook into our application.

There are two example applications in this repo to take a look at:

 * [tool](https://github.com/progrium/go-extpoints/tree/master/examples/tool) ([extpoint docs](http://godoc.org/github.com/progrium/go-extpoints/examples/tool/extpoints)), a command line tool with subcommands and lifecycle hooks as extension points
 * [daemon](https://github.com/progrium/go-extpoints/tree/master/examples/daemon), ... doesn't exist yet

 ## Extension Point Meta API

 All interfaces defined in your `extpoints` subpackage will be turned into global singleton extension point objects, using the pluralized name of the interface. These extension points implement this simple API:

 ```
 	type ExtensionPoint interface {
 		Register(extension <Interface>) bool
 		RegisterNamed(extension <Interface>, name string) bool
 		Lookup(name string) (<Interface>, bool)
 		All() map[string]<Interface>
 	}
 ```

 The `extpoints` subpackage will also have top-level registration functions that will run objects through all known extension points, registering with any that are based on an interface the object implements. They return string names of the interfaces they were registered against.
 
 ```
 	func Register(extension interface{}) []string
 	func RegisterNamed(extension interface{}, name string) []string
 ```

 ## License

 BSD