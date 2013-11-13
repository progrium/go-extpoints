package plugins

import(
	"io/ioutil"
	"strings"
	"os"
	"log"
	"path/filepath"
	"reflect"
	"github.com/robertkrimen/otto"
)

var pluginMap map[string]*otto.Otto
var implementations map[string]map[string]struct{}

func init() {
	pluginMap = make(map[string]*otto.Otto)
	implementations = make(map[string]map[string]struct{})
}

func loadRuntime(name, source string) (*otto.Otto, error) {
	runtime := otto.New()
	runtime.Set("implements", func(call otto.FunctionCall) otto.Value {
		interface_ := call.Argument(0).String()
		// todo: locks
		if implementations[name] == nil {
			implementations[name] = make(map[string]struct{})
		}
		implementations[name][interface_] = struct{}{}
    	return otto.UndefinedValue()
	})
	runtime.Run(source)
	return runtime, nil
}

func LoadString(name, source string) error {
	runtime, err := loadRuntime(name, source)
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
	return LoadString(strings.Split(filepath.Base(path), ".")[0], string(data))
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

func Register(ptr interface{}) {
	v := reflect.ValueOf(ptr).Elem()
	t := reflect.TypeOf(ptr).Elem()
	v2 := reflect.New(t).Elem()
	
	pluginField := v2.FieldByName("Plugin")
	typePtr := reflect.New(pluginField.Type().Out(0)).Interface()
	pluginFunc := func(params []reflect.Value) []reflect.Value {
		plugin := getPlugin(params[0].String(), typePtr)
		return []reflect.Value{ reflect.ValueOf(plugin) }
	}
	pluginField.Set(reflect.MakeFunc(pluginField.Type(), pluginFunc))
	
	pluginsField := v2.FieldByName("Plugins")
	pluginsFunc := func(params []reflect.Value) []reflect.Value {
		s := reflect.MakeSlice(reflect.SliceOf(pluginField.Type().Out(0)), 0, 0)
		for _, p := range getPlugins(typePtr) {
			s = reflect.Append(s, reflect.ValueOf(p))
		}
		return []reflect.Value{ s }
	}
	pluginsField.Set(reflect.MakeFunc(pluginsField.Type(), pluginsFunc))
	
	v.Set(v2)
}


func getPlugins(ptr interface{}) []interface{} {
	var plugins []interface{}
	for plugin, _ := range pluginMap {
		p := getPlugin(plugin, ptr)
		if p != nil {
			plugins = append(plugins, p)
		}
	}
	return plugins
}

func getPlugin(name string, ptr interface{}) interface{} {
	t := reflect.TypeOf(ptr).Elem()
	if !pluginImplements(name, t.Name()) {
		return nil
	}
	v := reflect.New(t).Elem()
	// loop over fields defined in type of ptr, 
	// replacing them in v with implementations
	for i, n := 0, t.NumField(); i < n; i++ {
		field := v.Field(i)
		structField := t.Field(i)
		newFunc := func(args []reflect.Value) []reflect.Value {
			if pluginMap[name] == nil {
				return []reflect.Value{ reflect.ValueOf(nil) }	
			}
        	value := callPlugin(name, structField.Name, convertArgs(args))
        	if value != nil {
        		return []reflect.Value{ reflect.ValueOf(value) }	
        	}
        	return []reflect.Value{}
    	}
		field.Set(reflect.MakeFunc(field.Type(), newFunc))
	}
	return v.Interface()
}

func pluginImplements(plugin, interface_ string) bool {
	_, found := implementations[plugin]
	if found {
		_, found = implementations[plugin][interface_]
		return found
	}
	return false
}

func callPlugin(plugin, function string, args []interface{}) interface{} {
	value, err := pluginMap[plugin].Call(function, nil, args...)
	if err != nil {
		log.Println("plugins:call: ", err)
		return nil
	}
	exported, _ := value.Export()
	return exported
}

func convertArgs(args []reflect.Value) []interface{} {
	var converted []interface{}
	for _, v := range args {
		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			converted = append(converted, v.Int())
		case reflect.String:
			converted = append(converted, v.String())
		default:
			log.Fatal("Unsupported type for argument: ", v.Type())
		}
	}
	return converted
}