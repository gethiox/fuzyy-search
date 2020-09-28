package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"fuzzy-search/internal/app/gutenbergsearch"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type serviceMock struct {
	errToReturn error
}

func (s *serviceMock) Search(title, phrase string) (string, error) {
	if s.errToReturn != nil {
		return "", s.errToReturn
	}
	return "", nil
}
func (s *serviceMock) Close() error {
	return nil
}

func (s *serviceMock) setErrToReturn(err error) {
	s.errToReturn = err
}

func testApp() (*httptest.Server, *serviceMock) {
	r := mux.NewRouter()
	searchService := &serviceMock{}
	r.Handle("/search", search(searchService))
	return httptest.NewServer(r), searchService
}

func Test_Search(t *testing.T) {
	type testCase struct {
		description        string
		payload            []byte
		expectedStatusCode int
	}

	var testCases = []testCase{
		{
			description:        "Unexpected format of payload (not JSON)",
			payload:            []byte(`asdf`),
			expectedStatusCode: http.StatusBadRequest,
		}, {
			description:        "Missing fields",
			payload:            []byte(`{"some": "payload"}`),
			expectedStatusCode: http.StatusBadRequest,
		}, {
			description:        "Missing 'phrase' field",
			payload:            []byte(`{"title": "some_title"}`),
			expectedStatusCode: http.StatusBadRequest,
		}, {
			description:        "Missing 'title' field",
			payload:            []byte(`{"phrase": "some_phrase"}`),
			expectedStatusCode: http.StatusBadRequest,
		}, {
			description:        "Payload OK",
			payload:            []byte(`{"title": "some_title", "phrase": "some_phrase"}`),
			expectedStatusCode: http.StatusOK,
		}, {
			description:        "Wrong `title` field type",
			payload:            []byte(`{"title": 10, "phrase": "some_phrase"}`),
			expectedStatusCode: http.StatusBadRequest,
		}, {
			description:        "Wrong `phrase` field type",
			payload:            []byte(`{"title": "some_title", "phrase": true}`),
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	ts, _ := testApp()
	defer ts.Close()

	client := http.Client{
		Timeout: time.Second * 2,
	}

	for i, tc := range testCases {
		name := fmt.Sprintf("%d:%s", i, tc.description)
		t.Run(name, func(t *testing.T) {
			request, err := http.NewRequest(http.MethodPost, ts.URL+"/search", bytes.NewBuffer(tc.payload))
			assert.Nil(t, err)
			res, err := client.Do(request)
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedStatusCode, res.StatusCode)
		})
	}
}

func Test_SearchError(t *testing.T) {
	type testCase struct {
		description        string
		errorReturned      error
		expectedStatusCode int
	}

	var testCases = []testCase{
		{
			description:        "Phrase not found",
			errorReturned:      gutenbergsearch.ErrPhraseNotFound,
			expectedStatusCode: http.StatusNotFound,
		}, {
			description:        "Processing too long",
			errorReturned:      gutenbergsearch.ErrTooLong,
			expectedStatusCode: http.StatusInternalServerError,
		},
	}

	ts, service := testApp()
	defer ts.Close()

	client := http.Client{
		Timeout: time.Second * 2,
	}

	for i, tc := range testCases {
		name := fmt.Sprintf("%d:%s", i, tc.description)
		t.Run(name, func(t *testing.T) {
			service.setErrToReturn(tc.errorReturned)
			payload := `{"title": "some_title", "phrase": "some_phrase"}`
			request, err := http.NewRequest(http.MethodPost, ts.URL+"/search", bytes.NewBuffer([]byte(payload)))
			assert.Nil(t, err)
			res, err := client.Do(request)
			assert.Nil(t, err)
			assert.Equal(t, tc.expectedStatusCode, res.StatusCode)
		})
	}
}
