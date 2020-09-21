package dataprovider

import (
	"errors"
	"regexp"
)

var queryResponseRegex = regexp.MustCompile(
	"(?s)" + // flag includes newline character into the scope of dot (.)
		"<li class=\"booklink\">.*?" +
		"<a class=\"link\" href=\"(?P<linkref>.*?)\".*?" +
		"<span class=\"title\">(?P<title>.*?)</span>.*?" +
		"<span class=\"subtitle\">(?P<subtitle>.*?)</span>.*?" +
		"</li>",
)

func parseBooks(responseBody string) ([]Book, error) {
	var books []Book

	matches := queryResponseRegex.FindAllStringSubmatch(responseBody, -1)
	if len(matches) == 0 {
		return books, errors.New("no results")
	}

	for _, match := range matches {
		book, err := NewBook(match[2], match[3], match[1])
		if err != nil {
			continue
		}
		books = append(books, book)
	}

	return books, nil
}
