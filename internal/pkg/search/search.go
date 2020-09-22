package search

type Result struct {
	Phrase     string // Matched phrase
	PosS, PosE int    // Position of phrase in book content
}

type Searcher interface {
	Search(content *string, phrase string) (Result, error)
}

func NewSearcher() Searcher {
	panic("TODO")
}
