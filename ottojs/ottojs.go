package ottojs

import (
	"github.com/robertkrimen/otto"
	"reflect"
	"strings"
)

type Runtime struct {
	plugins map[string]*otto.Otto
}

func GetRuntime() *Runtime {
	return &Runtime{
		plugins: make(map[string]*otto.Otto),
	}
}

func (r Runtime) FileExtension() string {
	return ".js"
}

func (r Runtime) InitPlugin(name, source string, implements func(string)) error {
	context := otto.New()
	context.Set("implements", func(call otto.FunctionCall) otto.Value {
		implements(call.Argument(0).String())
		return otto.UndefinedValue()
	})
	context.Run("implements_ = implements")
	context.Run(source)
	r.plugins[name] = context
	return nil
}

func (r Runtime) CallPlugin(name, function string, args []interface{}) (interface{}, error) {
	value, err := r.plugins[name].Call(function, nil, args...)
	if err != nil {
		return nil, err
	}
	exported, _ := value.Export()
	return exported, nil
}

func (r Runtime) SetGlobals(globals map[string]interface{}) {
	for k, v := range globals {
		for pluginName := range r.plugins {
			context := r.plugins[pluginName]
			if reflect.TypeOf(v).Kind() == reflect.Func {
				setValueAtPath(context, k, funcToOtto(context, reflect.ValueOf(v)))
			} else {
				setValueAtPath(context, k, v)
			}
		}
	}
}

func setValueAtPath(context *otto.Otto, path string, value interface{}) {
	parts := strings.Split(path, ".")
	parentCount := len(parts) - 1
	if parentCount > 0 {
		parentPath := strings.Join(parts[0:parentCount], ".")
		parent, err := context.Object("(" + parentPath + ")")
		if err != nil {
			emptyObject, _ := context.Object(`({})`)
			setValueAtPath(context, parentPath, emptyObject)
		}
		parent, _ = context.Object("(" + parentPath + ")")
		parent.Set(parts[parentCount], value)
	} else {
		context.Set(path, value)
	}
}

func funcToOtto(context *otto.Otto, fn reflect.Value) interface{} {
	return func(call otto.FunctionCall) otto.Value {
		convertedArgs := make([]reflect.Value, 0)
		for _, v := range call.ArgumentList {
			exported, _ := v.Export()
			convertedArgs = append(convertedArgs, reflect.ValueOf(exported))
		}
		ret := fn.Call(convertedArgs)
		if len(ret) > 0 {
			val, _ := context.ToValue(ret[0].Interface())
			return val
		} else {
			return otto.UndefinedValue()
		}
	}
}
