package plugins

import (
	"./ottojs"
	"testing"
)

type TestPlugin struct {
	Run func() string
}

var TestPluginExt struct {
	Plugin  func(string) TestPlugin
	Plugins func() []TestPlugin
}

func Test_GlobalFuncReturnStr(t *testing.T) {
	gloabls := map[string]interface{}{
		"TestReturnStr": func() string {
			return "testing 123"
		},
	}

	test_runtime := ottojs.GetRuntime()
	RegisterRuntime(test_runtime)

	js_test_plugin := `implements("TestPlugin")

function Run() {
    return typeof(TestReturnStr())
}
`
	LoadString("test_return_string", js_test_plugin, test_runtime)
	SetGlobals(gloabls)
	ExtensionPoint(&TestPluginExt)

	test_plugin := TestPluginExt.Plugin("test_return_string")
	output := test_plugin.Run()
	if output != "string" {
		t.Errorf("Got: \"%s\", Expected: \"%s\"", output, "string")
	}
}

func Test_GlobalFuncReturnInt(t *testing.T) {
	gloabls := map[string]interface{}{
		"TestReturnInt": func() int {
			return 1
		},
	}

	test_runtime := ottojs.GetRuntime()
	RegisterRuntime(test_runtime)

	js_test_plugin := `implements("TestPlugin")

function Run() {
    return typeof(TestReturnInt())
}
`
	LoadString("test_return_int", js_test_plugin, test_runtime)
	SetGlobals(gloabls)
	ExtensionPoint(&TestPluginExt)

	test_plugin := TestPluginExt.Plugin("test_return_int")
	output := test_plugin.Run()
	if output != "number" {
		t.Errorf("Got: \"%s\", Expected: \"%s\"", output, "number")
	}
}

func Test_GlobalFuncReturnFloat(t *testing.T) {
	gloabls := map[string]interface{}{
		"TestReturnFloat": func() float64 {
			return 1.0
		},
	}

	test_runtime := ottojs.GetRuntime()
	RegisterRuntime(test_runtime)

	js_test_plugin := `implements("TestPlugin")

function Run() {
    return typeof(TestReturnFloat())
}
`
	LoadString("test_return_float", js_test_plugin, test_runtime)
	SetGlobals(gloabls)
	ExtensionPoint(&TestPluginExt)

	test_plugin := TestPluginExt.Plugin("test_return_float")
	output := test_plugin.Run()
	if output != "number" {
		t.Errorf("Got: \"%s\", Expected: \"%s\"", output, "number")
	}
}

func Test_GlobalFuncReturnArray(t *testing.T) {
	gloabls := map[string]interface{}{
		"TestReturnArray": func() []int {
			arr := []int{1, 2, 3}
			return arr
		},
	}

	test_runtime := ottojs.GetRuntime()
	RegisterRuntime(test_runtime)

	js_test_plugin := `implements("TestPlugin")

function Run() {
    return typeof(TestReturnArray())
}
`
	LoadString("test_return_array", js_test_plugin, test_runtime)
	SetGlobals(gloabls)
	ExtensionPoint(&TestPluginExt)

	test_plugin := TestPluginExt.Plugin("test_return_array")
	output := test_plugin.Run()
	if output != "object" {
		t.Errorf("Got: \"%s\", Expected: \"%s\"", output, "object")
	}
}

func Test_GlobalFuncReturnMap(t *testing.T) {
	gloabls := map[string]interface{}{
		"TestReturnMap": func() map[string]interface{} {
			m := map[string]interface{}{
				"str": "testing",
				"int": 123,
			}
			return m
		},
	}

	test_runtime := ottojs.GetRuntime()
	RegisterRuntime(test_runtime)

	test_return_map_1 := `implements("TestPlugin")

function Run() {
    return typeof(TestReturnMap())
}
`
	LoadString("test_return_map_1", test_return_map_1, test_runtime)

	test_return_map_2 := `implements("TestPlugin")

function Run() {
    return typeof(TestReturnMap()["str"])
}
`
	LoadString("test_return_map_2", test_return_map_2, test_runtime)

	test_return_map_3 := `implements("TestPlugin")

function Run() {
    return typeof(TestReturnMap()["int"])
}
`
	LoadString("test_return_map_3", test_return_map_3, test_runtime)

	SetGlobals(gloabls)
	ExtensionPoint(&TestPluginExt)

	//test_return_map_1
	test_plugin_1 := TestPluginExt.Plugin("test_return_map_1")
	output := test_plugin_1.Run()
	if output != "object" {
		t.Errorf("Got: \"%s\", Expected: \"%s\"", output, "object")
	}

	//test_return_map_2
	test_plugin_2 := TestPluginExt.Plugin("test_return_map_2")
	output = test_plugin_2.Run()
	if output != "string" {
		t.Errorf("Got: \"%s\", Expected: \"%s\"", output, "string")
	}

	//test_return_map_3
	test_plugin_3 := TestPluginExt.Plugin("test_return_map_3")
	output = test_plugin_3.Run()
	if output != "number" {
		t.Errorf("Got: \"%s\", Expected: \"%s\"", output, "number")
	}
}
