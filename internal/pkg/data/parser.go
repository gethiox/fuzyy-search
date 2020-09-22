package data

import (
	"errors"
	"regexp"
)

var (
	findBooksRegexStage1 = regexp.MustCompile("(?s)<li class=\"booklink\">.*?</li>")
	findBooksRegexStage2 = regexp.MustCompile(
		"(?s)" + // flag includes newline character into the scope of dot (.)
			"<li class=\"booklink\">.*?" +
			"<a class=\"link\" href=\"(?P<linkref>.*?)\".*?" + // group 1
			"<span class=\"title\">(?P<title>.*?)</span>.*?" + // group 2
			"<span class=\"subtitle\">(?P<subtitle>.*?)</span>.*?" + // group 3
			"</li>",
	)
	findTxtLinkrefRegexStage1 = regexp.MustCompile("(?s)<a href.+?</a>")
	findTxtLinkrefRegexStage2 = regexp.MustCompile(
		"(?s)" +
			"<a href=\"(?P<txtlinkref>\\S*?)\".*?" + // group 1
			"title=\"Download\">Plain Text.*?" +
			"</a>",
	)
)

// findBooks parses book search query
func findBooks(responseBody string) ([]Book, error) {
	var books []Book

	matches := findBooksRegexStage1.FindAllString(responseBody, -1)
	if len(matches) == 0 {
		return books, errors.New("no results")
	}

	for _, stage2 := range matches {
		match := findBooksRegexStage2.FindStringSubmatch(stage2)
		if match == nil {
			continue
		}

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

// findTxtLinkref parses book entry page
func findTxtLinkref(responseBody string) (string, error) {
	matches := findTxtLinkrefRegexStage1.FindAllString(responseBody, -1)
	if len(matches) == 0 {
		return "", errors.New("no results")
	}

	for _, stage2 := range matches {
		matches := findTxtLinkrefRegexStage2.FindStringSubmatch(stage2)
		if matches == nil {
			continue
		}
		return matches[1], nil
	}

	return "", errors.New("no results")
}
