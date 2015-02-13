package testextpoints

type Dummy interface {
	Do()
}

func init() {
	dummies.Register(new(dummyA), "")
}

type dummyA struct{}

func (da *dummyA) Do() {
}
