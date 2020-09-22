package data

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Book struct {
	Title  string
	Author string

	bookLinkref string // eg. "/ebooks/34505", can be treated as unique ID
}

type errNoLinkRef struct{}

func (e *errNoLinkRef) Error() string {
	return "txt linkref is not available"
}

func errTxtLinkRefNotAvailable() error { return &errNoLinkRef{} }

var ErrTxtLinkRefNotAvailable = errTxtLinkRefNotAvailable()

func (b *Book) ID() string {
	return b.bookLinkref
}

func NewBook(title, author, linkref string) (Book, error) {
	if linkref == "" {
		return Book{}, errors.New("linkref is required")
	}

	if !strings.HasPrefix(linkref, "/") {
		linkref = "/" + linkref
	}

	return Book{
		Title:       title,
		Author:      author,
		bookLinkref: linkref,
	}, nil
}

type Provider interface {
	GetBooks(title string) ([]Book, error)
	DownloadBook(book Book) (string, error)
}

type httpProvider struct {
	Client http.Client

	domain    string
	userAgent string
}

func (p *httpProvider) baseUrl() string {
	return fmt.Sprintf("https://%s", p.domain)
}

// GetBooks return search results of title query.
func (p *httpProvider) GetBooks(title string) ([]Book, error) {
	// Function could return more results than first 25 popular matches by iterating through pages
	// but I decided not to du it, searching in 25 books for a query should be more than enough already
	query := url.QueryEscape(title)

	requestUrl := p.baseUrl() + "/ebooks/search/?query=" + query + "&submit_search=Go%21"

	request, err := http.NewRequest(http.MethodGet, requestUrl, nil)
	if err != nil {
		return []Book{}, fmt.Errorf("preparing request failed: %w", err)
	}
	request.Header.Set("User-Agent", p.userAgent)

	resp, err := p.Client.Do(request)
	if err != nil {
		return []Book{}, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []Book{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []Book{}, err
	}

	books, err := findBooks(string(body))
	if err != nil {
		return books, fmt.Errorf("books not found: %w", err)
	}

	return books, nil
}

// findTxtLinkRef returns linkref to txt version of given book, returns empty string if txt version is not available
func (p *httpProvider) findTxtLinkRef(book Book) (string, error) {
	requestUrl := fmt.Sprintf("%s%s", p.baseUrl(), book.bookLinkref)

	request, err := http.NewRequest(http.MethodGet, requestUrl, nil)
	if err != nil {
		return "", fmt.Errorf("preparing request failed: %w", err)
	}

	request.Header.Set("User-Agent", p.userAgent)

	resp, err := p.Client.Do(request)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response failed: %w", err)
	}

	linkref, err := findTxtLinkref(string(body))
	if err != nil {
		return "", fmt.Errorf("finding txt linkref failed: %w", err)
	}

	return linkref, nil
}

// DownloadBook tries to download text version of given book entry.
func (p *httpProvider) DownloadBook(book Book) (string, error) {
	linkRef, err := p.findTxtLinkRef(book)
	if err != nil {
		return "", fmt.Errorf("failed to get txt linkref: %w", err)
	}
	// some sleep to pretend real-human operation
	time.Sleep(time.Second * 2) // TODO: randomize

	requestUrl := fmt.Sprintf("%s%s", p.baseUrl(), linkRef)

	request, err := http.NewRequest(http.MethodGet, requestUrl, nil)
	if err != nil {
		return "", fmt.Errorf("preparing request failed: %w", err)
	}

	request.Header.Set("User-Agent", p.userAgent)

	resp, err := p.Client.Do(request)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response failed: %w", err)
	}

	return string(body), nil
}

func NewProvider() Provider {
	return &httpProvider{
		Client: http.Client{
			Timeout: time.Second * 60,
		},
		domain:    "www.gutenberg.org",
		userAgent: "Mozilla/5.0 (X11; Linux x86_64; rv:80.0) Gecko/20100101 Firefox/80.0",
	}
}
