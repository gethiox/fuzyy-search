package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"fuzzy-search/internal/app/gutenbergsearch"
	"fuzzy-search/internal/pkg/context"
	"fuzzy-search/internal/pkg/data"
	search2 "fuzzy-search/internal/pkg/search"

	"github.com/gorilla/mux"
)

type Payload struct {
	Title  *string `json:"title"`
	Phrase *string `json:"phrase"`
}

type ErrorMessage struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func newError(code, message string) []byte {
	msg := ErrorMessage{
		Code:    code,
		Message: message,
	}
	data, _ := json.Marshal(msg)
	return data
}

const (
	ErrMissingFiled   = "missing_filed"
	ErrJSONParse      = "bad_payload"
	ErrServerError    = "request_failed"
	ErrPhraseNotFound = "phrase_not_found"
)

func search(searchService gutenbergsearch.Searcher) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var payload Payload

		rawData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("read request failed: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write(newError(ErrJSONParse, "failed to read request"))
			return
		}

		err = json.Unmarshal(rawData, &payload)
		if err != nil {
			log.Printf("unmarshall failed: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			message := fmt.Sprintf("incorrect payload format: %s", err)
			_, _ = w.Write(newError(ErrJSONParse, message))
			return
		}

		if payload.Phrase == nil || payload.Title == nil {
			var missingFields []string

			if payload.Title == nil {
				missingFields = append(missingFields, "'title'")
			}
			if payload.Phrase == nil {
				missingFields = append(missingFields, "'phrase'")
			}

			w.WriteHeader(http.StatusBadRequest)
			message := fmt.Sprintf("missing fields: %s", strings.Join(missingFields, ", "))
			_, _ = w.Write(newError(ErrMissingFiled, message))
			return
		}

		var emptyFields []string
		for field, value := range map[string]string{
			"title":  *payload.Title,
			"phrase": *payload.Phrase,
		} {
			if value == "" {
				emptyFields = append(emptyFields, field)
			}
		}
		if len(emptyFields) > 0 {
			fields := strings.Join(emptyFields, ", ")
			w.WriteHeader(http.StatusBadRequest)
			message := fmt.Sprintf("fields cannot be empty: %s", fields)
			_, _ = w.Write(newError(ErrMissingFiled, message))
			return
		}

		// Search() could receive request context for processing cancellation purpose
		result, err := searchService.Search(*payload.Title, *payload.Phrase)
		if err != nil {
			switch {
			case errors.Is(err, gutenbergsearch.ErrPhraseNotFound):
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write(newError(ErrPhraseNotFound, "given phrase not found in books that matches given title"))
				return
			case errors.Is(err, gutenbergsearch.ErrTooLong):
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write(newError(ErrServerError, "requested processing exceeded allowed time"))
				return
			}

			log.Printf("Unexpected error ocurred: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write(newError(ErrServerError, "Something blows up on backend side, check logs for more details"))
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(result))
		return
	})
}

func prepareSearchService(cfg *Config) gutenbergsearch.Searcher {
	answerCache := gutenbergsearch.NewCache(cfg.answerCache, time.Hour*4, time.Minute*31)
	listingCache := gutenbergsearch.NewCache(cfg.listingCache, time.Hour*4, time.Minute*10)
	contentCache := gutenbergsearch.NewCache(cfg.contentCache, time.Hour, time.Minute*10)

	return gutenbergsearch.NewSearcher(
		8,
		answerCache,
		listingCache,
		contentCache,
		data.NewProvider(cfg.providerUserAgent, cfg.providerTimeout),
		context.NewProvider(),
		search2.NewSearcher(cfg.searchMaxDistance, cfg.searchRandomResult),
		[2]time.Duration{cfg.downloadDelayMin, cfg.downloadDelayMax},
	)
}

func main() {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	cfg := GetConfig()
	log.Printf("Loaded config:")
	log.Printf("%#v", cfg)

	searchService := prepareSearchService(cfg)
	defer func() {
		err := searchService.Close()
		if err != nil {
			log.Printf("Error ocurred during service shutdown: %s", err)
		}
	}()

	router := mux.NewRouter()
	router.Handle("/search", search(searchService))

	srv := &http.Server{
		Handler:      router,
		Addr:         ":8000",
		WriteTimeout: cfg.serverWriteTimeout,
		ReadTimeout:  cfg.serverReadTimeout,
	}

	log.Printf("Starting http server")

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			log.Printf("Unexpected server error: %s", err)
		}
		log.Print("Server stopped")
	}()

	<-done
	log.Print("Closing application...")
	err := srv.Close()
	if err != nil {
		log.Printf("Error occurred during close of webserver: %s", err)
	}
	log.Print("Application closed")
}
