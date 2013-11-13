package plugins

import(
	"io/ioutil"
	"strings"
	"os"
	"log"
	"path/filepath"
	"reflect"
	"errors"
)

// todo: locks!
var pluginMap map[string]Runtime
var implementations map[string]map[string]struct{}
var runtimes map[*Runtime]struct{}

type Runtime interface {
	FileExtension() string
	InitPlugin(name, source string, implements func(string)) error
	CallPlugin(name, function string, args []interface{}) (interface{}, error)
	LoadEnvironment(environment interface{})
	ConvertArgs(args []reflect.Value) []interface{}
}

func init() {
	pluginMap = make(map[string]Runtime)
	implementations = make(map[string]map[string]struct{})
	runtimes = make(map[*Runtime]struct{})
}

func RegisterRuntime(runtime Runtime) {
	runtimes[&runtime] = struct{}{}
}

func RegisterEnvironment(environ interface{}) {
	// TODO
}

func StaticPlugin(plugin *interface{}, interfaces []string) {
	// TODO
}

func LoadString(name, source string, runtime Runtime) error {
	err := runtime.InitPlugin(name, source, func(interfaceName string) {
		// todo: locks
		if implementations[name] == nil {
			implementations[name] = make(map[string]struct{})
		}
		implementations[name][interfaceName] = struct{}{}
	})
	if err != nil {
		return err
	}
	pluginMap[name] = runtime
	return nil
}

func LoadFile(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	var runtime Runtime
	for r := range runtimes {
		runtime = *r
		if strings.HasSuffix(path, runtime.FileExtension()) {
			break
		}
	}
	if runtime == nil {
		return errors.New("plugins: no runtime found to handle: "+path)
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
		if strings.HasSuffix(entry.Name(), ".js") {
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
		return []reflect.Value{ reflect.ValueOf(plugin) }
	}
	pluginField.Set(reflect.MakeFunc(pluginField.Type(), pluginFunc))
	
	pluginsField := extImpl.FieldByName("Plugins")
	pluginsFunc := func(params []reflect.Value) []reflect.Value {
		s := reflect.MakeSlice(reflect.SliceOf(pluginField.Type().Out(0)), 0, 0)
		for _, p := range getPluginProxies(extInterface) {
			s = reflect.Append(s, reflect.ValueOf(p))
		}
		return []reflect.Value{ s }
	}
	pluginsField.Set(reflect.MakeFunc(pluginsField.Type(), pluginsFunc))
	
	extPtr.Set(extImpl)
}


func getPluginProxies(extInterface interface{}) []interface{} {
	var plugins []interface{}
	for plugin, _ := range pluginMap {
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
			runtime := pluginMap[name]
			if runtime == nil {
				return []reflect.Value{ reflect.ValueOf(nil) }	
			}
			value, err := runtime.CallPlugin(name, structField.Name, runtime.ConvertArgs(args))
			if err != nil {
				log.Println("plugins:", err)
			}
        	if value != nil {
        		return []reflect.Value{ reflect.ValueOf(value) }	
        	}
        	return []reflect.Value{}
    	}
		field.Set(reflect.MakeFunc(field.Type(), newFunc))
	}
	return pluginProxy.Interface()
}

func hasImplementation(pluginName, interfaceName string) bool {
	_, found := implementations[pluginName]
	if found {
		_, found = implementations[pluginName][interfaceName]
		return found
	}
	return false
}
