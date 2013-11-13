# go-plugins

*early and experimental, but already badass*

Let users extend your Go applications with JavaScript and eventually Lua (or anything).

Thanks to [Otto](https://github.com/robertkrimen/otto), [Go reflection](http://golang.org/pkg/reflect/), and inspired by [Trac component architecture](http://trac.edgewall.org/wiki/TracDev/ComponentArchitecture).

## Using plugins

First, you define an "extension point". This involves writing an interface that plugins can implement, and then an extension point stub. Here is a simple observer pattern extension point:

	type ProgramObserver struct {
		ProgramStarted func()
		ProgramEnded func()
	}
	var ProgramObserverExt struct {
		Plugin func(string) ProgramObserver
		Plugins func() []ProgramObserver
	}

Now in the main() of your program, you load any plugins and register the extension point:

	plugins.LoadFromPath()
	plugins.Register(&ProgramObserverExt)

Now use the extension point in your program. `.Plugins()` gets you all plugins implementing that extension point interface, whereas `.Plugin(name)` lets you get a specific plugin by name.:

	for _, observer := range ProgramObserverExt.Plugins() {
		observer.ProgramStarted()
	}

	fmt.Println("Hello World")

	for _, observer := range ProgramObserverExt.Plugins() {
		observer.ProgramEnded()
	}

When we run the output is pretty boring:

	Hello World

Now let's write a plugin. We'll call it `happy.js`:

	implements("ProgramObserver")

	function ProgramStarted() {
		console.log("Yay! It's starting!")
	}

	function ProgramEnded() {
		console.log("Yay! It's over!")
	}

A plugin can implement any number of extension point interfaces. Now when we run our program? 

	Yay! It's starting!
	Hello World
	Yay! It's over!

Change the text in the plugin and run again. No need to recompile your Go. Add another plugin. Remove all plugins. It just works. Now go look at [all the examples](https://github.com/progrium/go-plugins/tree/master/examples).