package context

import (
	"errors"
)

type Provider interface {
	ProvideContext(content string, PosS, PosB int) (string, error)
}

type provider struct{}

func (p *provider) ProvideContext(content string, PosS, PosB int) (string, error) {
	var lastChar rune
	for i, c := range content[PosS:] {
		if c == '\n' && lastChar == '\n' {
			return content[PosS : PosS+i-1], nil
		}
		lastChar = c
	}
	return "", errors.New("context not selected")
}

func NewProvider() Provider {
	return &provider{}
}
