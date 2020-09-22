package gutenbergsearch

import (
	"fmt"
	"log"
	"time"

	"fuzzy-search/internal/pkg/context"
	"fuzzy-search/internal/pkg/data"
	"fuzzy-search/internal/pkg/search"

	"github.com/patrickmn/go-cache"
)

type book struct {
	title, author string
	uniqueID      string
	content       *string
}

type result struct {
	book   book
	result string
}

type Searcher interface {
	Search(title, phrase string) (string, error)
}

type searcher struct {
	bookCache, queryCache *cache.Cache

	dataProvider        data.Provider
	contextProvider     context.Provider
	searchEngine        search.Searcher
	searchEngineWorkers int // per request
}

func NewSearcher() Searcher {
	return &searcher{
		bookCache:           cache.New(time.Hour, time.Minute*10),
		queryCache:          cache.New(time.Hour*12, time.Minute*31),
		dataProvider:        data.NewProvider(),
		contextProvider:     context.NewProvider(),
		searchEngine:        search.NewSearcher(),
		searchEngineWorkers: 8,
	}
}

// twoPartCacheKey generate unique key for cache usage
func twoPartCacheKey(a, b string) string {
	// TODO: a and b could be validated against possible edge-case scenarios
	return a + "/" + b
}

func (s *searcher) Search(title, phrase string) (string, error) {
	cachedAnswer, ok := s.queryCache.Get(twoPartCacheKey(title, phrase))
	if ok {
		log.Println("found cached query result")
		return cachedAnswer.(string), nil
	}

	log.Printf("Searching books with \"%s\" title", title)

	bookPositions, err := s.dataProvider.GetBooks(title)
	if err != nil {
		return "", fmt.Errorf("searching books failed: %w", err)
	}

	log.Printf("Found %d books", len(bookPositions))

	if len(bookPositions) < 1 {
		return "", fmt.Errorf("no books available for this title")
	}

	var booksToAnalyze = make(chan book, 25)
	defer close(booksToAnalyze)
	var resultChan = make(chan result, 1)
	defer close(resultChan)

	// Run search engine workers
	for i := 0; i < s.searchEngineWorkers; i++ {
		go func() {
			bookZeroValue := book{}
			for {
				book := <-booksToAnalyze
				if book == bookZeroValue {
					return // channel closed
				}
				searcgResult, err := s.searchEngine.Search(book.content, phrase)
				if err != nil {
					log.Printf("no result for \"%s\" (%s) book: %s", book.title, book.author, err)
					continue
				}

				withContext, err := s.contextProvider.ProvideContext(searcgResult.Phrase, searcgResult.PosS, searcgResult.PosE)
				if err != nil {
					log.Printf("failed to provide context for \"%s\" match: %s", searcgResult.Phrase, err)
					continue
				}

				s.queryCache.Set(twoPartCacheKey(title, phrase), withContext, cache.DefaultExpiration)
				resultChan <- result{
					book:   book,
					result: withContext,
				}
			}
		}()
	}

	var booksProcessedFromCache = make(map[string]bool) // key: book unique ID

	// Gather all currently cached books
	for _, bookPosition := range bookPositions {
		cachedBook, ok := s.bookCache.Get(bookPosition.ID())
		if !ok {
			continue
		}

		booksProcessedFromCache[bookPosition.ID()] = true
		log.Println("found cached book")
		bookContent := cachedBook.(string)

		booksToAnalyze <- book{
			title:   bookPosition.Title,
			author:  bookPosition.Author,
			content: &bookContent, // TODO: not sure that's a good idea
		}
	}

	var downloadQueue = make(chan data.Book, 25)
	defer close(downloadQueue)

	// Missing books downloader
	go func() {
		zeroValue := data.Book{}
		for {
			bookPosition := <-downloadQueue
			if bookPosition == zeroValue {
				return // closed channel
			}
			for i := 0; i < 3; i++ { // 3 download tries
				content, err := s.dataProvider.DownloadBook(bookPosition)
				if err != nil {
					time.Sleep(time.Second * 5)
				}
				s.bookCache.Set(bookPosition.ID(), content, cache.DefaultExpiration)
				booksToAnalyze <- book{
					title:   bookPosition.Title,
					author:  bookPosition.Author,
					content: &content, // TODO: not sure that's a good idea
				}
			}
		}
	}()

	// Queue up missing books to download
	go func() {
		for _, bookPosition := range bookPositions {
			_, ok := booksProcessedFromCache[bookPosition.ID()]
			if ok {
				continue
			}

			downloadQueue <- bookPosition
		}
	}()

	result := <-resultChan
	s.queryCache.Set(twoPartCacheKey(title, phrase), result.result, cache.DefaultExpiration)
	return result.result, nil
}
