package main

import (
	"testing"

	"github.com/progrium/go-extpoints/tests/extpoints"
)

var noops = extpoints.Noops
var noopFactories = extpoints.NoopFactories
var transformers = extpoints.StringTransformers

func TestLookupSuccess(t *testing.T) {
	_, ok := noops.Lookup("noop")
	if !ok {
		t.Fatal("Lookup returned not ok for registered extension")
	}
}

func TestLookupFail(t *testing.T) {
	_, ok := noops.Lookup("yesop")
	if ok {
		t.Fatal("Lookup returned ok for non-existent extension")
	}
}

func TestUsingComponent(t *testing.T) {
	upper, _ := transformers.Lookup("upper")
	if upper.Transform("string") != "STRING" {
		t.Fatal("Used component, but didn't work as expected")
	}
}

func TestUsingFuncComponent(t *testing.T) {
	factory, _ := noopFactories.Lookup("noopFactory")
	n := factory()
	if _, ok := n.(extpoints.Noop); !ok {
		t.Fatal("Used component, but didn't work as expected")
	}
}
