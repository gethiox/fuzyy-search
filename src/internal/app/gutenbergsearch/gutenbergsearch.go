package gutenbergsearch

import (
	context2 "context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"fuzzy-search/internal/pkg/context"
	"fuzzy-search/internal/pkg/data"
	"fuzzy-search/internal/pkg/search"

	"github.com/patrickmn/go-cache"
)

var (
	ErrPhraseNotFound = errors.New("phrase not found")
)

type book struct {
	title, author string
	uniqueID      string
	content       string
}

type result struct {
	book   book
	result string
}

type Searcher interface {
	Search(title, phrase string) (string, error)
}

func randomRange(a, b time.Duration) time.Duration {
	rawI := rand.Intn(int(b - a))
	return time.Duration(int(a) + rawI)
}

type searcher struct {
	bookCache, bookPositionsCache, queryCache *cache.Cache

	dataProvider        data.Provider
	contextProvider     context.Provider
	searchEngine        search.Searcher
	downloadDelay       [2]time.Duration // min, max
	searchEngineWorkers int              // per request
}

func NewSearcher() Searcher {
	rand.Seed(time.Now().UnixNano())

	return &searcher{
		bookCache:          cache.New(time.Hour, time.Minute*10), // heaviest cache, stores downloaded book content
		bookPositionsCache: cache.New(time.Hour*4, time.Minute*10),
		queryCache:         cache.New(time.Hour*4, time.Minute*31),

		dataProvider:        data.NewProvider(),
		contextProvider:     context.NewProvider(),
		searchEngine:        search.NewSearcher(4, true),
		downloadDelay:       [2]time.Duration{time.Second, time.Second * 2},
		searchEngineWorkers: 8,
	}
}

// twoPartCacheKey generate unique key for cache usage
func twoPartCacheKey(a, b string) string {
	// TODO: a and b could be validated against possible edge-case scenarios
	return a + "/" + b
}

func (s *searcher) getBookPositions(title string) ([]data.Book, error) {
	cachedBookPositions, ok := s.bookPositionsCache.Get(title)
	if ok {
		bookPositions := cachedBookPositions.([]data.Book)
		log.Printf("Read %d book positions from cache", len(bookPositions))
		return bookPositions, nil
	}

	bookPositions, err := s.dataProvider.GetBooks(title)
	if err != nil {
		return bookPositions, fmt.Errorf("downloading book positions failed: %w", err)
	}
	log.Printf("Read %d book positions from external source", len(bookPositions))
	s.bookPositionsCache.Set(title, bookPositions, cache.DefaultExpiration)
	return bookPositions, nil
}

func (s *searcher) Search(title, phrase string) (string, error) {
	var exit bool

	cachedAnswer, ok := s.queryCache.Get(twoPartCacheKey(title, phrase))
	if ok {
		log.Println("found cached query result")
		return cachedAnswer.(string), nil
	}

	log.Printf("Searching books with \"%s\" title", title)
	bookPositions, err := s.getBookPositions(title)
	if err != nil {
		return "", fmt.Errorf("getBookPositions failed: %w", err)
	}

	if len(bookPositions) < 1 {
		return "", fmt.Errorf("no books available for this title")
	}

	resultCtx, resultCancel := context2.WithCancel(context2.Background())
	jobsCtx, jobsCancel := context2.WithCancel(context2.Background())

	var booksToAnalyze = make(chan book, 25)
	defer close(booksToAnalyze)

	var resultChan = make(chan result, 1)
	defer close(resultChan)

	jobs := sync.WaitGroup{}
	for i := 0; i < len(bookPositions); i++ {
		jobs.Add(1)
	}

	log.Printf("Run job monitor")
	go func() {
		log.Printf("Waiting to process all jobs")
		jobs.Wait()
		log.Printf("All jobs processed")
		jobsCancel()
	}()

	log.Printf("Run search engine workers")

	// Run search engine workers
	for i := 0; i < s.searchEngineWorkers; i++ {
		go searchWorker(i, s, resultCtx, title, phrase, booksToAnalyze, resultChan, &jobs)
	}

	var booksProcessedFromCache = make(map[string]bool) // key: book unique ID

	// Gather all currently cached books
	for _, bookPosition := range bookPositions {
		cachedBook, ok := s.bookCache.Get(bookPosition.ID())
		if !ok {
			continue
		}

		booksProcessedFromCache[bookPosition.ID()] = true
		log.Printf("Load book from cache (\"%s\" - %s)", bookPosition.Title, bookPosition.Author)
		bookContent := cachedBook.(string)

		if exit {
			break
		}
		booksToAnalyze <- book{
			title:    bookPosition.Title,
			author:   bookPosition.Author,
			uniqueID: bookPosition.ID(),
			content:  bookContent,
		}
	}

	var downloadQueue = make(chan data.Book, 25)
	defer close(downloadQueue)

	// Missing books downloader
	go downloadWorker(s, resultCtx, downloadQueue, booksToAnalyze)

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

	select {
	case result := <-resultChan:
		resultCancel()
		s.queryCache.Set(twoPartCacheKey(title, phrase), result.result, cache.DefaultExpiration)
		log.Printf("result found")
		return result.result, nil
	case <-jobsCtx.Done():
		resultCancel()
		log.Printf("result not found")
		return "", ErrPhraseNotFound
	}
}

func searchWorker(n int, s *searcher, resultCtx context2.Context, title, phrase string, booksToAnalyze <-chan book, resultChan chan<- result, jobs *sync.WaitGroup) {
	bookZeroValue := book{}
	log.Printf("search worker %d started", n)
	for {
		book := <-booksToAnalyze
		if book == bookZeroValue {
			log.Printf("Search worker %d closed", n)
			return // channel closed
		}
		searchResult, err := s.searchEngine.Search(book.content, phrase)
		if err != nil {
			log.Printf("no result for book (\"%s\" - %s [%s]): %s", book.title, book.author, book.uniqueID, err)
			jobs.Done()
			continue
		}

		withContext, err := s.contextProvider.ProvideContext(book.content, searchResult.PosS, searchResult.PosE)
		if err != nil {
			log.Printf("failed to provide context for \"%s\" match: %s", searchResult.Phrase, err)
			jobs.Done()
			continue
		}

		s.queryCache.Set(twoPartCacheKey(title, phrase), withContext, cache.DefaultExpiration)
		select {
		case <-resultCtx.Done():
			log.Printf("Search worker %d closed", n)
			return
		case <-time.After(time.Millisecond * 5):
			resultChan <- result{
				book:   book,
				result: withContext,
			}
			jobs.Done()
			break
		}
	}
}

func downloadWorker(s *searcher, resultCtx context2.Context, downloadQueue <-chan data.Book, booksToAnalyze chan<- book) {
	zeroValue := data.Book{}
main:
	for {
		bookPosition := <-downloadQueue
		if bookPosition == zeroValue {
			return // closed channel
		}
		downloadTries := 3
		for i := 0; i < downloadTries; i++ {
			time.Sleep(randomRange(s.downloadDelay[0], s.downloadDelay[1]))

			content, err := s.dataProvider.DownloadBook(bookPosition)
			if err != nil {
				if errors.Is(err, data.ErrTxtLinkRefNotAvailable) {
					// book position does not include text version
					log.Printf(
						"text book not available (\"%s\" - %s [%s]): %s",
						bookPosition.Title, bookPosition.Author, bookPosition.ID(),
						err,
					)
					// return an empty book anyway for caching purpose, TODO: could it be handled better
					booksToAnalyze <- book{
						title:    bookPosition.Title,
						author:   bookPosition.Author,
						uniqueID: bookPosition.ID(),
						content:  "",
					}
					s.bookCache.Set(bookPosition.ID(), "", cache.DefaultExpiration)
					continue main

				}

				log.Printf("download error: %s", err)
				continue main
			}
			s.bookCache.Set(bookPosition.ID(), content, cache.DefaultExpiration)
			select {
			case <-resultCtx.Done():
				log.Printf("Downloading interrupted")
				return
			case <-time.After(time.Millisecond * 5):
				booksToAnalyze <- book{
					title:    bookPosition.Title,
					author:   bookPosition.Author,
					uniqueID: bookPosition.ID(),
					content:  content,
				}
				break
			}
		}
	}
}
