package search

import (
	"errors"
	"strings"
)

type Result struct {
	Phrase     string // Matched phrase
	PosS, PosE int    // Position of phrase in book content
}

type Searcher interface {
	Search(content *string, phrase string) (Result, error)
}

type localSearcher struct{}

func (l *localSearcher) Search(content *string, phrase string) (Result, error) {
	idx := strings.Index(*content, phrase)
	if idx == -1 {
		return Result{}, errors.New("pattern not found")
	}

	idx2 := idx + len(phrase)

	return Result{
		Phrase: (*content)[idx:idx2],
		PosS:   idx,
		PosE:   idx2,
	}, nil
}

func NewSearcher() Searcher {
	return &localSearcher{}
}
