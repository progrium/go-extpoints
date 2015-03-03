package extpoints

type Noop interface {
	Noop()
}

type StringTransformer interface {
	Transform(input string) string
}

type NoopFactory func() Noop
