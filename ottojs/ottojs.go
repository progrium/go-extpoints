package ottojs

import (
	"github.com/robertkrimen/otto"
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

func (r Runtime) LoadEnvironment(environment interface{}) {

}
