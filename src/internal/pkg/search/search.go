package search

import (
	"errors"
	"math/rand"
	"strings"
	"time"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

type Result struct {
	Phrase     string // Matched phrase
	PosS, PosE int    // Position of phrase in book content
}

type Searcher interface {
	Search(content *string, phrase string) (Result, error)
}

type localSearcher struct {
	maxDistance  int
	randomResult bool
}

type Indexes struct {
	a, b int
}

func BookFields(s string, cap int) ([]string, []Indexes) {
	fields := make([]string, 0, cap)
	indexes := make([]Indexes, 0, cap)

	if len(s) < 1 {
		return fields, indexes
	}

	var state int
	// determining initial state
	switch s[0] {
	case ' ', '\n', '\t':
		state = WhitespaceState
	default:
		state = TextState
	}

	var textBeginning int
	for i, r := range s {
		switch r {
		case ' ', '\n', '\t', '\r':
			if state == TextState {
				fields = append(fields, s[textBeginning:i])
				indexes = append(indexes, Indexes{textBeginning, i})
				state = WhitespaceState
			}
		default:
			if state == WhitespaceState {
				state = TextState
				textBeginning = i
			}
		}
	}
	return fields, indexes
}

const (
	WhitespaceState = iota
	TextState
)

func (l *localSearcher) Search(content *string, phrase string) (Result, error) {
	var indexes = make([]Indexes, 0)

	contentFields, contentIndexes := BookFields(*content, 100)
	phraseFields := strings.Fields(phrase)

	if len(phraseFields) < 1 {
		return Result{}, errors.New("pattern not found")
	}

	firstPhraseField := phraseFields[0]
	ranks := fuzzy.RankFindFold(firstPhraseField, contentFields)

	var filteredRanks fuzzy.Ranks
	for _, rank := range ranks {
		if rank.Distance > l.maxDistance {
			continue
		}
		filteredRanks = append(filteredRanks, rank)
	}
	if len(filteredRanks) == 0 {
		return Result{}, errors.New("pattern not found")
	}

	for _, rank := range filteredRanks {
		for i, phraseField := range phraseFields[1:] {
			targetIndex := rank.OriginalIndex + 1 + i
			secondRanks := fuzzy.RankFindFold(phraseField, []string{contentFields[targetIndex]})
			if len(secondRanks) == 0 {
				// no match
				continue
			}
			if secondRanks[0].Distance > l.maxDistance {
				continue
			}
			if i == len(phraseFields)-2 {
				// achieved full query pass successfully
				firstWord := rank.OriginalIndex
				lastWord := rank.OriginalIndex + len(phraseFields) - 1
				indexes = append(indexes, Indexes{contentIndexes[firstWord].a, contentIndexes[lastWord].b})
			}
		}
	}

	if len(indexes) == 0 {
		return Result{}, errors.New("pattern not found")
	}

	var choice Indexes
	if l.randomResult {
		choice = indexes[rand.Intn(len(indexes))]
	} else {
		choice = indexes[0]
	}

	return Result{
		Phrase: (*content)[choice.a:choice.b],
		PosS:   choice.a,
		PosE:   choice.b,
	}, nil
}

func NewSearcher(maxDistance int, randomResult bool) Searcher {
	if randomResult {
		rand.Seed(time.Now().UnixNano())
	}

	return &localSearcher{
		maxDistance:  maxDistance,
		randomResult: true,
	}
}
