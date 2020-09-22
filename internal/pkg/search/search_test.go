package search

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalSearcher_SearchSimpleHappyPath(t *testing.T) {
	type testCase struct {
		Description    string
		Searcher       Searcher
		Content        string
		Phrase         string
		expectedResult Result
	}

	var testCases = []testCase{
		{
			Description:    "no fuzzing, exact match",
			Searcher:       NewSearcher(0, 0, 0),
			Content:        "abcd",
			Phrase:         "bc",
			expectedResult: Result{Phrase: "bc", PosS: 1, PosE: 3},
		}, {
			Description:    "fuzzing, insertion 1",
			Searcher:       NewSearcher(1, 0, 0),
			Content:        "abcde",
			Phrase:         "bd",
			expectedResult: Result{Phrase: "bcd", PosS: 1, PosE: 4},
		}, {
			Description:    "fuzzing, replacement 1",
			Searcher:       NewSearcher(0, 1, 0),
			Content:        "abcde",
			Phrase:         "bXd",
			expectedResult: Result{Phrase: "bcd", PosS: 1, PosE: 4},
		}, {
			Description:    "fuzzing, deletion 1",
			Searcher:       NewSearcher(0, 0, 1),
			Content:        "abcde",
			Phrase:         "bcXd",
			expectedResult: Result{Phrase: "bcd", PosS: 1, PosE: 4},
		},
	}

	for i, tc := range testCases {
		name := fmt.Sprintf("%d:%s", i, tc.Description)
		t.Run(name, func(t *testing.T) {
			result, err := tc.Searcher.Search(&tc.Content, tc.Phrase)
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestLocalSearcher_SearchFuzzingInsertion(t *testing.T) {
	type testCase struct {
		Description    string
		Searcher       Searcher
		Content        string
		Phrase         string
		expectedResult Result
	}

	var testCases = []testCase{
		{
			Description:    "insertion 1",
			Searcher:       NewSearcher(1, 0, 0),
			Content:        "some content to search",
			Phrase:         "contentto",
			expectedResult: Result{Phrase: "content to", PosS: 5, PosE: 15},
		}, {
			Description:    "insertion 2, two chars in series",
			Searcher:       NewSearcher(2, 0, 0),
			Content:        "some content to search",
			Phrase:         "contento",
			expectedResult: Result{Phrase: "content to", PosS: 5, PosE: 15},
		}, {
			Description:    "insertion 2, two chars in different places",
			Searcher:       NewSearcher(2, 0, 0),
			Content:        "some content to search",
			Phrase:         "cntentto",
			expectedResult: Result{Phrase: "content to", PosS: 5, PosE: 15},
		},
	}

	for i, tc := range testCases {
		name := fmt.Sprintf("%d:%s", i, tc.Description)
		t.Run(name, func(t *testing.T) {
			result, err := tc.Searcher.Search(&tc.Content, tc.Phrase)
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}
func TestLocalSearcher_SearchFuzzingReplacement(t *testing.T) {
	type testCase struct {
		Description    string
		Searcher       Searcher
		Content        string
		Phrase         string
		expectedResult Result
	}

	var testCases = []testCase{
		{
			Description:    "replacement 1",
			Searcher:       NewSearcher(0, 1, 0),
			Content:        "some content to search",
			Phrase:         "content tX",
			expectedResult: Result{Phrase: "content to", PosS: 5, PosE: 15},
		}, {
			Description:    "replacement 2, two chars in series",
			Searcher:       NewSearcher(0, 2, 0),
			Content:        "some content to search",
			Phrase:         "content XX",
			expectedResult: Result{Phrase: "content to", PosS: 5, PosE: 15},
		}, {
			Description:    "replacement 2, two chars in different places",
			Searcher:       NewSearcher(0, 2, 0),
			Content:        "some content to search",
			Phrase:         "conXent oY",
			expectedResult: Result{Phrase: "content to", PosS: 5, PosE: 15},
		},
	}

	for i, tc := range testCases {
		name := fmt.Sprintf("%d:%s", i, tc.Description)
		t.Run(name, func(t *testing.T) {
			result, err := tc.Searcher.Search(&tc.Content, tc.Phrase)
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}
func TestLocalSearcher_SearchFuzzingDeletion(t *testing.T) {
	type testCase struct {
		Description    string
		Searcher       Searcher
		Content        string
		Phrase         string
		expectedResult Result
	}

	var testCases = []testCase{
		{
			Description:    "deletion 1",
			Searcher:       NewSearcher(0, 0, 1),
			Content:        "some content to search",
			Phrase:         "content Wto",
			expectedResult: Result{Phrase: "content to", PosS: 5, PosE: 15},
		}, {
			Description:    "deletion 2, two chars in series",
			Searcher:       NewSearcher(0, 0, 2),
			Content:        "some content to search",
			Phrase:         "content YWto",
			expectedResult: Result{Phrase: "content to", PosS: 5, PosE: 15},
		}, {
			Description:    "deletion 2, two chars in different places",
			Searcher:       NewSearcher(0, 0, 2),
			Content:        "some content to search",
			Phrase:         "contXent Yto",
			expectedResult: Result{Phrase: "content to", PosS: 5, PosE: 15},
		},
	}

	for i, tc := range testCases {
		name := fmt.Sprintf("%d:%s", i, tc.Description)
		t.Run(name, func(t *testing.T) {
			result, err := tc.Searcher.Search(&tc.Content, tc.Phrase)
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}

func TestLocalSearcher_SearchFuzzingInsertionReplacementDeletion(t *testing.T) {
	type testCase struct {
		Description    string
		Searcher       Searcher
		Content        string
		Phrase         string
		expectedResult Result
	}

	var testCases = []testCase{
		{
			Description:    "insertion 1, replacement 1, deletion 1",
			Searcher:       NewSearcher(1, 1, 1),
			Content:        "some content to search",
			Phrase:         "cotenR Dto",
			expectedResult: Result{Phrase: "content to", PosS: 5, PosE: 15},
		}, {
			Description:    "insertion 2, replacement 2, deletion 2",
			Searcher:       NewSearcher(2, 2, 2),
			Content:        "some content to search",
			Phrase:         "RtRnDt tDo",
			expectedResult: Result{Phrase: "content to", PosS: 5, PosE: 15},
		},
	}

	for i, tc := range testCases {
		name := fmt.Sprintf("%d:%s", i, tc.Description)
		t.Run(name, func(t *testing.T) {
			result, err := tc.Searcher.Search(&tc.Content, tc.Phrase)
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}
