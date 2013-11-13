# go-plugins (experimental)

Let users extend your Go applications with JavaScript and eventually Lua.

## Using plugins

First, you define an "extension point". This is just an interface that plugins can implement and more or less a factory thing. Here is a simple observer pattern:

	type ProgramObserver struct {
		ProgramStarted func()
		ProgramEnded func()
	}
	var ProgramObserverExt struct {
		Plugin func(string) ProgramObserver
		Plugins func() []ProgramObserver
	}

When you run your program, you load any plugins and register the extension point:

	plugins.LoadFromPath()
	plugins.Register(&ProgramObserverExt)

Now use the extension point in your program:

	for _, observer := range ProgramObserverExt.Plugins() {
		observer.ProgramStarted()
	}

	fmt.Println("Hello World")

	for _, observer := range ProgramObserverExt.Plugins() {
		observer.ProgramEnded()
	}

Now let's write a plugin. We'll call it `happy.js`:

	implements("ProgramObserver")

	function ProgramStarted() {
		console.log("Yay! It's starting!")
	}

	function ProgramEnded() {
		console.log("Yay! It's over!")
	}

When we run our program? 

	Yay! It's starting!
	Hello World
	Yay! It's over!

Change the text in the plugin and run again. No need to recompile your Go. Add another plugin. Remove all plugins. It just works. Now go look at [all the examples](https://github.com/progrium/go-plugins/tree/master/examples).