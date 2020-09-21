package search

type Book struct {
	Title  string
	Author string

	Content *string
}

type Result struct {
	Phrase     string // Matched phrase
	PosS, PosE int    // Position of phrase in book content
}

type Searcher interface {
	Search(book Book, phrase string) (Result, error)
}

func NewSearcher() Searcher {
	panic("TODO")
}
