package context

type Provider interface {
	ProvideContext(content *string, PosS, PosB int) (string, error)
}

func NewProvider() Provider {
	panic("TODO")
}
