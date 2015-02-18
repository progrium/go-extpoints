package main

func init() {
	dummies.Register(new(dummyA), "")
}

type dummyA struct{}

func (da *dummyA) Do() {
}
