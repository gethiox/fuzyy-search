package gutenbergsearch

import (
	context2 "context"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"sync"
	"time"

	"fuzzy-search/internal/pkg/context"
	"fuzzy-search/internal/pkg/data"
	"fuzzy-search/internal/pkg/search"
)

var (
	ErrPhraseNotFound = errors.New("phrase not found")
	ErrTooLong        = errors.New("request took too long")
)

type downloadJobs struct {
	ctx           context2.Context
	downloadQueue <-chan data.Book
	outputQueue   chan<- book
}

type searchJobs struct {
	ctx         context2.Context
	phrase      string
	searchQueue <-chan book
	outputQueue chan<- result
}

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
	io.Closer
}

func randomDurationRange(a, b time.Duration) time.Duration {
	rawI := rand.Intn(int(b - a))
	return time.Duration(int(a) + rawI)
}

// twoPartCacheKey generate unique key for cache usage
func twoPartCacheKey(a, b string) string {
	// TODO: a and b could be validated against possible edge-case scenarios
	return a + "/" + b
}

type searcher struct {
	answerCache, listingCache, contentCache Cache

	dataProvider        data.Provider
	contextProvider     context.Provider
	searchEngine        search.Searcher
	downloadDelay       [2]time.Duration // min, max
	searchEngineWorkers int              // per request

	tasksWg      sync.WaitGroup
	exit         chan bool
	downloadJobs chan downloadJobs
	searchJobs   chan searchJobs
}

// randomizedResult - return random search match if exist in the scope of currently processed book, can't work
// with AnswerCache enabled
func NewSearcher(
	searchWorkers int,
	answerCache, listingCache, contentCache Cache,
	dataProvider data.Provider,
	contextProvider context.Provider,
	searchEngine search.Searcher,
	downloadDelay [2]time.Duration, // min/max

) Searcher {
	rand.Seed(time.Now().UnixNano())

	s := &searcher{
		answerCache:  answerCache,
		listingCache: listingCache,
		contentCache: contentCache,

		dataProvider:        dataProvider,
		contextProvider:     contextProvider,
		searchEngine:        searchEngine,
		downloadDelay:       downloadDelay,
		searchEngineWorkers: searchWorkers,

		tasksWg:      sync.WaitGroup{},
		exit:         make(chan bool, 1),
		downloadJobs: make(chan downloadJobs, 5),
		searchJobs:   make(chan searchJobs, searchWorkers),
	}
	s.StartBackgroundTasks()
	return s
}

func (s *searcher) Close() error {
	close(s.exit)
	s.tasksWg.Wait()
	return nil
}

func (s *searcher) StartBackgroundTasks() {
	go s.downloadTask()
	s.tasksWg.Add(1)
	go s.searchTask()
	s.tasksWg.Add(1)
}

func (s *searcher) getBookPositions(title string) ([]data.Book, error) {
	cachedBookPositions, ok := s.listingCache.Get(title)
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
	s.listingCache.Set(title, bookPositions)
	return bookPositions, nil
}

func (s *searcher) Search(title, phrase string) (string, error) {
	var exit bool

	cachedAnswer, ok := s.answerCache.Get(twoPartCacheKey(title, phrase))
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

	var booksToAnalyze = make(chan book, 25)
	// downloadTask will close this channel

	var resultChan = make(chan result, 1)
	// searchTask will close this channel

	var booksProcessedFromCache = make(map[string]bool) // key: book unique ID

	// Gather all currently cached books
	for _, bookPosition := range bookPositions {
		cachedBook, ok := s.contentCache.Get(bookPosition.ID())
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

	downloadQueue := make(chan data.Book, 25)

	downloadCtx, downloadCancel := context2.WithCancel(context2.Background())
	defer downloadCancel()

	s.downloadJobs <- downloadJobs{
		ctx:           downloadCtx,
		downloadQueue: downloadQueue,
		outputQueue:   booksToAnalyze,
	}

	searchCtx, searchCancel := context2.WithCancel(context2.Background())
	defer searchCancel()

	s.searchJobs <- searchJobs{
		ctx:         searchCtx,
		phrase:      phrase,
		searchQueue: booksToAnalyze,
		outputQueue: resultChan,
	}

	// Queue up missing books to download
	go func() {
		var scheduled int
		for _, bookPosition := range bookPositions {
			_, ok := booksProcessedFromCache[bookPosition.ID()]
			if ok {
				continue
			}
			downloadQueue <- bookPosition
			scheduled += 1
		}
		log.Printf("Scheduled %d books to download", scheduled)
	}()

	select {
	case result, ok := <-resultChan:
		if ok {
			s.answerCache.Set(twoPartCacheKey(title, phrase), result.result)
			log.Printf("result found! ('%s' - %s)", result.book.title, result.book.author)
			return result.result, nil
		}
		// processing ended but no result pushed on channel
		return "", ErrPhraseNotFound
	case <-time.After(time.Second * 120):
		// processing took too long
		return "", ErrTooLong
	}
}

// downloadTask constantly monitor incoming downloadJobs queue and starts goroutine for every incoming job
func (s *searcher) downloadTask() {
	log.Print("[[ downloadTask running ]]")
root:
	for {
		select {
		case <-s.exit:
			break root
		case job := <-s.downloadJobs:
			go func(job downloadJobs) {
				log.Print("[downloadTask] Download job acquired")
				defer close(job.outputQueue)
				startTime := time.Now()
			main:
				for bookToDownload := range job.downloadQueue {
					select {
					case <-job.ctx.Done():
						log.Printf("[DWorker] Downloading interrupted")
						return
					case <-time.After(time.Millisecond * 10):
					}

					var content string
					var err error
					var failedTries int

					downloadTries := 3

				try:
					for i := 0; i < downloadTries; i++ {

						sleepTime := randomDurationRange(s.downloadDelay[0], s.downloadDelay[1])
						log.Printf("[DWorker] sleeping for %s", sleepTime)
						time.Sleep(sleepTime)

						content, err = s.dataProvider.DownloadBook(bookToDownload)
						if err != nil {
							if errors.Is(err, data.ErrTxtLinkRefNotAvailable) {
								// this book position apparently does not include text version
								log.Printf(
									"[DWorker] Text book not available (\"%s\" - %s [%s]): %s",
									bookToDownload.Title, bookToDownload.Author, bookToDownload.ID(),
									err,
								)
								// preparing an empty content as successful download for caching purpose
								content = ""
								break try
							}

							log.Printf("[DWorker] Download error: %s (try %d/%d)", err, i+1, downloadTries)
							failedTries += 1
							continue try
						}
						// successful download
						break try
					}
					if failedTries == downloadTries {
						log.Print("[DWorker] Download failed: retries exceeded")
						continue main
					}

					endTime := time.Now()

					downloadedBook := book{
						title:    bookToDownload.Title,
						author:   bookToDownload.Author,
						uniqueID: bookToDownload.ID(),
						content:  content,
					}

					log.Printf("[DWorker] Book ('%s' - %s) downloaded in %s",
						bookToDownload.Title,
						bookToDownload.Author,
						endTime.Sub(startTime),
					)

					s.contentCache.Set(bookToDownload.ID(), content)

					select {
					case <-job.ctx.Done():
						log.Printf("[DWorker] Downloading interrupted")
						return
					case job.outputQueue <- downloadedBook:
						continue main
					case <-time.After(time.Second * 10):
						log.Printf("[DWorker] Data push timeout")
						continue main
					}
				}
				log.Print("[downloadTask] Download job finished")
			}(job)
		}
	}
	log.Print("[[ downloadTask closed ]]")
	s.tasksWg.Done()
}

// searchTask constantly monitor incoming searchJobs queue and starts goroutine for every incoming job
func (s *searcher) searchTask() {
	log.Print("[[ searchTask running ]]")
root:
	for {
		select {
		case <-s.exit:
			break root
		case job := <-s.searchJobs:
			log.Print("[searchTask] Search job acquired")
			workerWg := sync.WaitGroup{}
			go func(job searchJobs) {
				defer close(job.outputQueue)
				for i := 0; i < s.searchEngineWorkers; i++ {
					workerWg.Add(1)
					go func(workerID int) {
						defer workerWg.Done()
						for book := range job.searchQueue {
							select {
							case <-job.ctx.Done():
								log.Printf("[SWorker %d] Search interrupted", workerID)
								return
							case <-time.After(time.Millisecond * 10):
							}

							searchResult, err := s.searchEngine.Search(book.content, job.phrase)
							if err != nil {
								log.Printf("[SWorker %d] no result for book (\"%s\" - %s [%s]): %s", workerID, book.title, book.author, book.uniqueID, err)
								return
							}

							withContext, err := s.contextProvider.ProvideContext(book.content, searchResult.PosS, searchResult.PosE)
							if err != nil {
								log.Printf("[SWorker %d] failed to provide context for \"%s\" match: %s", workerID, searchResult.Phrase, err)
								return
							}

							s.answerCache.Set(twoPartCacheKey(book.title, job.phrase), withContext)
							r := result{
								book:   book,
								result: withContext,
							}

							select {
							case <-job.ctx.Done():
								log.Printf("[SWorker %d] Search interrupted", workerID)
								return
							case job.outputQueue <- r:
								return
							case <-time.After(time.Second * 1):
								log.Printf("[SWorker %d] Push search output timeout", workerID)
								return
							}
						}
					}(i)
				}

				log.Printf("[searchTask] %d search workers running", s.searchEngineWorkers)
				workerWg.Wait()
				log.Printf("[searchTask] search workers finishes")
			}(job)
		}
	}
	log.Print("[[ searchTask closed ]]")
	s.tasksWg.Done()
}
