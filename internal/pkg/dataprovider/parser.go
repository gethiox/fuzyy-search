package dataprovider

import (
	"errors"
	"regexp"
)

var queryResponseRegex = regexp.MustCompile(
	"(?s)" + // flag includes newline character into the scope of dot (.)
		"<li class=\"booklink\">.*?" +
		"<a class=\"link\" href=\"(?P<linkref>.*?)\".*?" + // 1
		"<span class=\"title\">(?P<title>.*?)</span>.*?" + // 2
		"<span class=\"subtitle\">(?P<subtitle>.*?)</span>.*?" + // 3
		"</li>",
)

// parseBooks parses book search query
func parseBooks(responseBody string) ([]Book, error) {
	var books []Book

	matches := queryResponseRegex.FindAllStringSubmatch(responseBody, -1)
	if len(matches) == 0 {
		return books, errors.New("no results")
	}

	for _, match := range matches {
		title := match[2]
		author := match[3]
		linkref := match[1]

		book, err := NewBook(title, author, linkref)
		if err != nil {
			continue
		}
		books = append(books, book)
	}

	return books, nil
}
