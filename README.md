# go-plugins

*early and experimental, but already badass*

Let users extend your Go applications with JavaScript or any other runtime you implement.

Thanks to [Otto](https://github.com/robertkrimen/otto), [Go reflection](http://golang.org/pkg/reflect/), and inspired by [Trac component architecture](http://trac.edgewall.org/wiki/TracDev/ComponentArchitecture).

## Using go-plugins

### Extension Points

First, you define an "extension point" for you application. This is basically a set of hooks that a plugin can implement to extend some aspect of your application. Defining an extension point involves writing an interface (though not a real interface) that plugins can implement, and then an extension point stub. Here is a simple observer pattern extension point:

	type ProgramObserver struct {
		ProgramStarted func()
		ProgramEnded func()
	}
	var ProgramObserverExt struct {
		Plugin func(string) ProgramObserver
		Plugins func() []ProgramObserver
	}

You also want to register your extension point for it to be active:

	plugins.Register(&ProgramObserverExt)

Now use the extension point in your program. `.Plugins()` gets you all plugins implementing that extension point interface, whereas `.Plugin(name)` lets you get a specific plugin by name. More often you use the former; the latter is used when you are using plugins to provide configurable backends. But here's `.Plugins()` in our app:

	for _, observer := range ProgramObserverExt.Plugins() {
		observer.ProgramStarted()
	}

	fmt.Println("Hello World")

	for _, observer := range ProgramObserverExt.Plugins() {
		observer.ProgramEnded()
	}

Without plugins loaded, when we run the output, it's pretty boring:

	Hello World

### Runtimes

Plugins run in runtimes, which define a scripting environment. Out of the box, we have `ottojs`, which is a JavaScript runtime based on Otto, a pure Go JavaScript interpreter. You can define your own runtimes to hook up and script in Python, Lua, or anything else. You can even support multiple runtimes at once! Just register them at the beginning of your program:

	plugins.RegisterRuntime(ottojs.GetRuntime())

### Plugins

We're about to write a plugin! We'll call it `happy.js` and put it in a `plugins` directory. This is the default place to look when you load with `plugins.LoadFromPath()`, which you can override with the `PLUGIN_PATH` environment variable. There are plenty of other ways to load plugins, but this is the easiest:

	plugins.LoadFromPath()

And now some JavaScript. Our `plugins/happy.js` file:

	implements("ProgramObserver")

	function ProgramStarted() {
		console.log("Yay! It's starting!")
	}

	function ProgramEnded() {
		console.log("Yay! It's over!")
	}

A plugin can implement any number of extension point interfaces by calling `implements()` multiple times. With this when we run our program:

	Yay! It's starting!
	Hello World
	Yay! It's over!

Change the text in the plugin and run again. No need to recompile your Go. Add another plugin. Remove all plugins. It just works. You can see the [full source for this example](https://github.com/progrium/go-plugins/tree/master/examples/simple) or look at [all the examples](https://github.com/progrium/go-plugins/tree/master/examples).