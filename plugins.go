package plugins

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

// todo: locks!
var plugin struct {
	loaded     map[string]runner
	interfaces map[string]map[string]struct{}
	runtimes   map[*Runtime]struct{}
	globals    map[string]interface{}
}

type runner interface {
	CallPlugin(name, function string, args []interface{}) (interface{}, error)
}

type Runtime interface {
	runner
	FileExtension() string
	InitPlugin(name, source string, implements func(string)) error
	SetGlobals(globals map[string]interface{})
}

func init() {
	plugin.loaded = make(map[string]runner)
	plugin.interfaces = make(map[string]map[string]struct{})
	plugin.runtimes = make(map[*Runtime]struct{})
	plugin.globals = make(map[string]interface{})
}

func RegisterRuntime(runtime Runtime) {
	plugin.runtimes[&runtime] = struct{}{}
}

func SetGlobals(globals map[string]interface{}) {
	for k, v := range globals {
		plugin.globals[k] = v
	}
	for name := range plugin.loaded {
		r, ok := plugin.loaded[name].(Runtime)
		if ok {
			r.SetGlobals(plugin.globals)
		}
	}
}

func GetGlobal(name string) interface{} {
	return plugin.globals[name]
}

func StaticPlugin(name string, instance interface{}, interfaces []string) {
	for _, interfaceName := range interfaces {
		registerInterface(name, interfaceName)
	}
	plugin.loaded[name] = staticRunner{instance}
}

func LoadString(name, source string, runtime Runtime) error {
	err := runtime.InitPlugin(name, source, func(interfaceName string) {
		registerInterface(name, interfaceName)
	})
	if err != nil {
		return err
	}
	plugin.loaded[name] = runtime.(runner)
	return nil
}

func findRuntimeForFile(path string) Runtime {
	var runtime Runtime
	for r := range plugin.runtimes {
		runtime = *r
		fileExt := runtime.FileExtension()
		if fileExt != "" && strings.HasSuffix(path, fileExt) {
			break
		}
	}
	return runtime
}

func LoadFile(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	runtime := findRuntimeForFile(path)
	if runtime == nil {
		return errors.New("plugins: no runtime found to handle: " + path)
	}
	name := strings.Split(filepath.Base(path), ".")[0]
	return LoadString(name, string(data), runtime)
}

func LoadFromPath() error {
	path := os.Getenv("PLUGIN_PATH")
	if path == "" {
		path = "plugins"
	}
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	for _, entry := range dir {
		runtime := findRuntimeForFile(path)
		if runtime != nil {
			err = LoadFile(path + "/" + entry.Name())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func ExtensionPoint(ext interface{}) {
	extPtr := reflect.ValueOf(ext).Elem()
	extImpl := reflect.New(reflect.TypeOf(ext).Elem()).Elem()

	pluginField := extImpl.FieldByName("Plugin")
	extInterface := reflect.New(pluginField.Type().Out(0)).Interface()
	pluginFunc := func(params []reflect.Value) []reflect.Value {
		plugin := getPluginProxy(params[0].String(), extInterface)
		return []reflect.Value{reflect.ValueOf(plugin)}
	}
	pluginField.Set(reflect.MakeFunc(pluginField.Type(), pluginFunc))

	pluginsField := extImpl.FieldByName("Plugins")
	pluginsFunc := func(params []reflect.Value) []reflect.Value {
		s := reflect.MakeSlice(reflect.SliceOf(pluginField.Type().Out(0)), 0, 0)
		for _, p := range getPluginProxies(extInterface) {
			s = reflect.Append(s, reflect.ValueOf(p))
		}
		return []reflect.Value{s}
	}
	pluginsField.Set(reflect.MakeFunc(pluginsField.Type(), pluginsFunc))

	extPtr.Set(extImpl)
}

func getPluginProxies(extInterface interface{}) []interface{} {
	var plugins []interface{}
	for plugin, _ := range plugin.loaded {
		p := getPluginProxy(plugin, extInterface)
		if p != nil {
			plugins = append(plugins, p)
		}
	}
	return plugins
}

func getPluginProxy(name string, extInterface interface{}) interface{} {
	extInterfaceType := reflect.TypeOf(extInterface).Elem()
	if !hasImplementation(name, extInterfaceType.Name()) {
		return nil
	}
	pluginProxy := reflect.New(extInterfaceType).Elem()
	// loop over fields defined in extInterfaceType,
	// replacing them in v with implementations
	for i, n := 0, extInterfaceType.NumField(); i < n; i++ {
		field := pluginProxy.Field(i)
		structField := extInterfaceType.Field(i)
		newFunc := func(args []reflect.Value) []reflect.Value {
			runner := plugin.loaded[name]
			if runner == nil {
				return []reflect.Value{reflect.ValueOf(nil)}
			}
			value, err := runner.CallPlugin(name, structField.Name, convertArgs(args))
			if err != nil {
				log.Println("plugins:", err)
			}
			if value != nil {
				return []reflect.Value{reflect.ValueOf(value)}
			}
			return []reflect.Value{}
		}
		field.Set(reflect.MakeFunc(field.Type(), newFunc))
	}
	return pluginProxy.Interface()
}

func hasImplementation(pluginName, interfaceName string) bool {
	_, found := plugin.interfaces[pluginName]
	if found {
		_, found = plugin.interfaces[pluginName][interfaceName]
		return found
	}
	return false
}

func registerInterface(pluginName, interfaceName string) {
	// todo: locks
	if plugin.interfaces[pluginName] == nil {
		plugin.interfaces[pluginName] = make(map[string]struct{})
	}
	plugin.interfaces[pluginName][interfaceName] = struct{}{}
}

func convertArgs(args []reflect.Value) []interface{} {
	var converted []interface{}
	for _, v := range args {
		converted = append(converted, exportValue(v))
	}
	return converted
}

func exportValue(value reflect.Value) interface{} {
	switch value.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return value.Int()
	case reflect.String:
		return value.String()
	default:
		log.Fatal("ottojs: Unsupported type for argument:", value.Type())
		return nil
	}
}

type staticRunner struct {
	plugin interface{}
}

func (r staticRunner) CallPlugin(name, function string, args []interface{}) (interface{}, error) {
	p := reflect.ValueOf(r.plugin).Elem()
	f := p.MethodByName(function)
	var argValues []reflect.Value
	for _, v := range args {
		argValues = append(argValues, reflect.ValueOf(v))
	}
	value := f.Call(argValues)
	if len(value) > 0 {
		return exportValue(value[0]), nil
	}
	return nil, nil
}
