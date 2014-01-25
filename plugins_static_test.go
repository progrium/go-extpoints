package plugins

import (
	"testing"
)

type TestStaticPlugin struct {
	Run func() string
	//Id string //Conflicts with plugins.go:184
	
}

var TestStaticPluginExt struct {
	Plugin func(string) TestStaticPlugin
	Plugins func() []TestStaticPlugin
}

type StaticPlugin1 struct {
	Id string
}

func (p StaticPlugin1) Run() string {
	return p.Id
}

type StaticPlugin2 struct {
	Id string
}

func (p StaticPlugin2) Run() string {
	return p.Id
}


func Test_MultiStatic(t *testing.T) {
        ExtensionPoint(&TestStaticPluginExt)
        StaticPlugin("1a", &StaticPlugin1{Id: "1a"}, []string{"TestStaticPlugin"})
        StaticPlugin("1b", &StaticPlugin1{Id: "1b"}, []string{"TestStaticPlugin"})
        StaticPlugin("2a", &StaticPlugin2{Id: "2a"}, []string{"TestStaticPlugin"})

        testPlugin1a := TestStaticPluginExt.Plugin("1a")
	output := testPlugin1a.Run()
	if output != "1a" {
		t.Errorf("Got: \"%s\", Expected: \"%s\"", output, "1a")
	}

        testPlugin1b := TestStaticPluginExt.Plugin("1b")
	output = testPlugin1b.Run()
	if output != "1b" {
		t.Errorf("Got: \"%s\", Expected: \"%s\"", output, "1b")
	}

        testPlugin2a := TestStaticPluginExt.Plugin("2a")
	output = testPlugin2a.Run()
	if output != "2a" {
		t.Errorf("Got: \"%s\", Expected: \"%s\"", output, "2a")
	}
}
