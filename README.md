# go-extpoints

This package, named short for "Go extension points", provides a simple component model for making your Go packages, libraries, and applications extensible in a standard way. 

It expands on the foundation of Go interfaces and provides a meta-API for accessing and registering "components", objects that implement one or more "extension point" interfaces. It not only lets third-party packages hook in as build-time extensions, but also encourages better organization of your own package and subpackages. 

## Getting the tool

	$ go install github.com/progrium/go-extpoints

## Quick Example

Here is a simple Go application that turns its `main()` into an extension point based on a `ProgramParticipant` interface. Here is our `main.go` that we'll say lives under `$GOPATH/src/github.com/quick/example`:

```
//go:generate go-extpoints
package main

import (
	"github.com/quick/example/extpoints"
)

func main() {
	for _, extension := range extpoints.ProgramParticipants.All() {
		extension.Main()
	}
}
```

We create an `extpoints` subpackage with a Go file in it to define our extension point interface. This is `extpoints/interfaces.go`:

```
package extpoints

type ProgramParticipant interface {
	Main()
}
```

This is where `go-extpoints` comes in. It hooks into `go generate` to produce extension point code in our `extpoints` subpackage around the interfaces defined there.

	$ go generate
	 ....
	$ go install

But it doesn't *do* anything! Let's make a builtin component that implements `ProgramParticipant`. Add a `builtin.go` file:

```
package main

import (
	"fmt"
	"github.com/quick/example/extpoints"
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
		Register(component <Interface>) bool
		RegisterNamed(component <Interface>, name string) bool
		Lookup(name string) (<Interface>, bool)
		All() map[string]<Interface>
	}
```

The `extpoints` subpackage will also have top-level registration functions that will run components through all known extension points, registering with any that are based on an interface the component implements. They return the names of the interfaces they were registered against.

```
	func Register(component interface{}) []string
	func RegisterNamed(component interface{}, name string) []string
```

## Why a subpackage?

There are number of reasons this turned out to be a very elegant solution. 

First, since we force the convention of a subpackage called `extpoints`, it makes it very easy to identify a package as having extension points from looking at the project tree. You then know where to look to find the interfaces that are exposed as extension points.

Second, it makes it clearer in your code when you're using extension points. You have to explicitly import the package, then call `extpoints.<ExtensionPoint>` when using them. This helps identify where extension points actually hook into your program.

Third, third-party packages have a well known package to import for registering. Whether you have extension points for a library package or a command with just a `main` package, there's always a definite `extpoints` package there to import.

Lastly, it produces its own GoDoc page. Extension points are designed to use existing documentation infrastructure. But in such a way that gives them their own namespace. Your extension point APIs are different than regular APIs. They're not APIs to call, but APIs to implement, specifically to extend your package. They're the "back office" APIs of your package.

## Usage Patterns

There are a number of API design patterns for extension point interfaces, however that's a deep subject. Instead, here are "mechanical" patterns of different ways to use the extension point meta-API:

#### Simple Iteration

#### Lookup Only One

#### Passing References

#### Match and Use

## Laying the Groundwork

Although this only seems to allow for compile-time extensibility, this is already quite a win. It means power users can build and compile in their own extensions that live outside your repository. 

However, it also lays the groundwork for other dynamic extensions. I've used this model to wrap extension points to implement components in embedded scripting languages, as hook scripts, or as remote plugin daemons via RPC. 

No matter how you're thinking about dynamic extensions later on, using `go-extpoints` gives you a lot of options. Once Go supports dynamic libraries? This will work perfectly with that too.

## Inspiration

This project and the model that it supports is a Go idiomatic port of the [component architecture](http://trac.edgewall.org/wiki/TracDev/ComponentArchitecture) used in Trac, which is written in Python. It's taken about a year to get this right in Go.

## License

BSD