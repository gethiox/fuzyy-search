package search

import (
	"io/ioutil"
	"os"
	"testing"
)

func loadTestBook(t *testing.T) Book {
	fd, err := os.OpenFile("search_test_book_content.txt", os.O_RDONLY, 0)
	if err != nil {
		t.Fatal("Failed to open book content file: ", err)
	}
	data, err := ioutil.ReadAll(fd)
	if err != nil {
		t.Fatal("Failed to load book content data: ", err)
	}

	content := string(data)
	return Book{
		Title:   "Title",
		Author:  "Author",
		Content: &content,
	}
}

func TestExpectedResult(t *testing.T) {
	book := loadTestBook(t)
	phrase := "oh romeo romeo"
	expectedResult := Result{}

	searcher := NewSearcher()

	result, err := searcher.Search(book, phrase)
	if err != nil {
		t.Fatal("Search failed; ", err)
	}

	if expectedResult.Phrase != result.Phrase {
		t.Fatalf("Phrase does not match, Wanted: \"%s\", Have: \"%s\"", expectedResult.Phrase, result.Phrase)
	}
	if expectedResult.PosS != result.PosS {
		t.Fatalf("PosS does not match, Wanted: \"%d\", Have: \"%d\"", expectedResult.PosS, result.PosS)
	}
	if expectedResult.PosE != result.PosE {
		t.Fatalf("PosE does not match, Wanted: \"%d\", Have: \"%d\"", expectedResult.PosE, result.PosE)
	}
}

func TestExpectedResultCorrectPosition(t *testing.T) {
	book := loadTestBook(t)
	phrase := "oh romeo romeo"
	expectedMatch := "O Romeo, Romeo"

	searcher := NewSearcher()

	result, err := searcher.Search(book, phrase)
	if err != nil {
		t.Fatal("Search failed; ", err)
	}

	match := (*book.Content)[result.PosS:result.PosE]

	if match != expectedMatch {
		t.Logf("Position range in result is incorrect (%d:%d)", result.PosS, result.PosE)
		t.Logf("Wanted: \"%s\"", expectedMatch)
		t.Logf("Have: \"%s\"", match)
	}
}
