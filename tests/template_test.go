package main

import (
	"testing"

	"./extpoints"
)

var dummies = extpoints.Dummies

func TestLookupSuccess(t *testing.T) {
	_, ok := dummies.Lookup("dummyA")
	if !ok {
		t.Fatal("Lookup returned not ok for registered extension")
	}
}

func TestLookupFail(t *testing.T) {
	_, ok := dummies.Lookup("dummyNotExist")
	if ok {
		t.Fatal("Lookup returned ok for non-existent extension")
	}
}
