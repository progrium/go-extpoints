# go-extpoints

This package, named short for "Go extension points", provides a simple component model for making your Go packages, libraries, and applications extensible in a standard way. 

It expands on the foundation of Go interfaces and provides a meta-API for accessing and registering "components", objects that implement one or more "extension point" interfaces. It not only lets third-party packages hook in as build-time extensions, but also encourages better organization of your own packages. 

## Getting the tool

	$ go install github.com/progrium/go-extpoints

## Quick Example

Here is a simple Go application that lets components hook into `main()` as subcommands by implementing an extension point interface called `Subcommand`. Assuming our package lives under `$GOPATH/src/github.com/quick/example`, here is our `main.go`:

```go
//go:generate go-extpoints
package main

import (
	"fmt"
	"os"
	
	"github.com/quick/example/extpoints"
)

var subcommands = extpoints.Subcommands

func usage() {
	fmt.Println("Available commands:\n")
	for name, _ := range subcommands.All() {
		fmt.Println(" - ", name)
	}
	os.Exit(2)
}

func main() {
	if len(os.Args) < 2 {
		usage()
	}
	cmd, exists := subcommands.Lookup(os.Args[1])
	if !exists {
		usage()
	}
	cmd.Run(os.Args[2:])
}
```

Note the extension point is referred to by the plural of our interface `Subcommand` and under the `extpoints` subpackage. We need to create that package with a Go file in it to define our extension point interface. This is `extpoints/interfaces.go`:

```go
package extpoints

type Subcommand interface {
	Run(args []string)
}
```

We use `go generate`, which calls `go-extpoints`, to produce extension point code in our `extpoints` subpackage around the interfaces defined there.

	$ go generate
	 ...
	$ go install

Okay, but it doesn't *do* anything! Let's make a builtin command component that implements `Subcommand`. Add a `hello.go` file:

```go
package main

import (
	"fmt"
	"github.com/quick/example/extpoints"
)

func init() {
	extpoints.Register(new(HelloComponent), "hello")
}

type HelloComponent struct {}

func (p *HelloComponent) Run(args []string) {
	fmt.Println("Hello world!")
}
```

Now when we build and run the app, it shows `hello` as a subcommand. 

Certainly, the value of extension points becomes clearer with larger applications and more interesting interfaces. But just consider now that the component defined in `hello.go` *could* exist in another package in another repo. You'd just have to import it and rebuild to let it hook into our application.

There are two more in-deptch example applications in this repo to take a look at:

 * [tool](https://github.com/progrium/go-extpoints/tree/master/examples/tool) ([extpoints](http://godoc.org/github.com/progrium/go-extpoints/examples/tool/extpoints)), a more realistic CLI tool with subcommands and lifecycle hooks
 * [daemon](https://github.com/progrium/go-extpoints/tree/master/examples/daemon), ... doesn't exist yet

## Extension Point Meta API

All interfaces defined in your `extpoints` subpackage will be turned into extension point singleton object variables, using the pluralized name of the interface. These extension point objects implement this simple meta-API:

```go
type <ExtensionPoint> interface {
	// if name is "", the component type is used
	Register(component <Interface>, name string) bool
	
	Unregister(name string) bool

	Lookup(name string) (<Interface>, bool)

	All() map[string]<Interface> // keyed by name
}
```

Your `extpoints` subpackage will also have top-level registration functions generated that will run components through all known extension points, registering or unregistering with any that are based on an interface the component implements. They return the names of the interfaces they were registered/unregistered with.

```
func Register(component interface{}, name string) []string
func Unregister(name string) []string
```

## Making it easy to install extensions

Assuming you tell third-party developers to call your `extpoints.Register` in their `init()`, you can link them with a side-effect import (using a blank import name). 

You can make this easy for users to enable/disable via comments, or add their own without worrying about messing with your code by having a separate `extensions.go` or `plugins.go` file with just these imports:

```go
package yourpackage

import (
	_ "github.com/you/some-extension"
	_ "github.com/third-party/another-extension"
)

```

Users can now just edit this file and `go build` or `go install`. 

## Usage Patterns

Here are different example ways to use extension points to interact with components:

#### Simple Iteration
```go
for _, listener := range extpoints.EventListeners.All() {
	listener.Notify(&MyEvent{})
}
```

#### Lookup Only One
```go
driverName := config.Get("storage-driver")
driver, registered := extpoints.StorageDrivers.Lookup(driverName)
if !registered {
	log.Fatalf("storage driver '%s' not installed", driverName)
}
driver.StoreObject(object)
```

#### Passing by Reference
```go
for _, filter := range extpoints.RequestFilters.All() {
	filter.FilterRequest(req)
}
```

#### Match and Use
```go
for _, handler := range extpoints.RequestHandlers.All() {
	if handler.MatchRequest(req) {
		handler.HandleRequest(req)
		break
	}
}
```

## Why the `extpoints` subpackage?

Since we force the convention of a subpackage called `extpoints`, it makes it very easy to identify a package as having extension points from looking at the project tree. You then know where to look to find the interfaces that are exposed as extension points.

Third-party packages have a well known package to import for registering. Whether you have extension points for a library package or a command with just a `main` package, there's always a definite `extpoints` package there to import.

It also makes it clearer in your code when you're using extension points. You have to explicitly import the package, then call `extpoints.<ExtensionPoint>` when using them. This helps identify where extension points actually hook into your program.

## Groundwork for Dynamic Extensions

Although this only seems to allow for compile-time extensibility, this itself is quite a win. It means power users can build and compile in their own extensions that live outside your repository. 

However, it also lays the groundwork for other dynamic extensions. I've used this model to wrap extension points for components in embedded scripting languages, as hook scripts, as remote plugin daemons via RPC, or all of the above implemented as components themselves! 

No matter how you're thinking about dynamic extensions later on, using `go-extpoints` gives you a lot of options. Once Go supports dynamic libraries? This will work perfectly with that, too.

## Inspiration

This project and component model is a lightweight, Go idiomatic port of the [component architecture](http://trac.edgewall.org/wiki/TracDev/ComponentArchitecture) used in Trac, which is written in Python. It's taken about a year to get this right in Go.

## License

BSD