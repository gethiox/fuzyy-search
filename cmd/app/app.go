package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"fuzzy-search/internal/app/gutenbergsearch"

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

var searchService gutenbergsearch.Searcher

func search(w http.ResponseWriter, r *http.Request) {
	// defer close(w)
	var payload Payload

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("read request failed: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(newError(ErrJSONParse, "failed to read request"))
		return
	}

	err = json.Unmarshal(data, &payload)
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

	if *payload.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(newError(ErrMissingFiled, "'title' field cannot be empty"))
		return
	}

	if *payload.Phrase == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write(newError(ErrMissingFiled, "'phrase' field cannot be empty"))
		return
	}

	result, err := searchService.Search(*payload.Title, *payload.Phrase)
	if err != nil {
		if errors.Is(err, gutenbergsearch.ErrPhraseNotFound) {
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write(newError(ErrPhraseNotFound, "given phrase not found in books that matches given title"))
			return
		}

		log.Printf("Unexpected error ocurred: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write(newError(ErrServerError, "Something blows up on backend side, check logs for more details"))
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(result))
}

func main() {
	searchService = gutenbergsearch.NewSearcher()

	router := mux.NewRouter()
	router.HandleFunc("/search", search)

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Starting http server")
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf("unexpected server error: %s", err)
	}
}
