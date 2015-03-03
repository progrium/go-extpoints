package main

import (
	"strings"

	"github.com/progrium/go-extpoints/tests/extpoints"
)

func init() {
	extpoints.Register(new(noop), "")                      // Noop
	extpoints.Register(new(uppercaseTransformer), "upper") // StringTransformer
	extpoints.NoopFactories.Register(noopFactory, "")
}

func noopFactory() extpoints.Noop {
	return new(noop)
}

type noop struct{}

func (n *noop) Noop() {
}

type uppercaseTransformer struct{}

func (t *uppercaseTransformer) Transform(input string) string {
	return strings.ToUpper(input)
}
