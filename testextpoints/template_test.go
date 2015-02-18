package testextpoints

import "testing"

var dummies = Dummies

func TestLookupSuccess(t *testing.T) {
	lookup(t, "dummyA", true)
}

func TestLookupFail(t *testing.T) {
	lookup(t, "dummyNotExist", false)
}

func lookup(t *testing.T, name string, expectedOk bool) {
	_, ok := dummies.Lookup(name)
	if ok != expectedOk {
		t.Errorf("Expected 'ok' flag %v, but got %v\n", expectedOk, ok)
	}
}
