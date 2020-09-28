## Build and run (Docker)

Make sure your Docker is installed and running, then simply run (in root repository directory):

```shell
./build/docker-build.sh && ./build/docker-run.sh
```

## Build and run (locally)

Make sure your Go is up to date and execute (in root repository directory):

```shell
./build/local-build.sh && ./build/local-run.sh
```

### Making requests

App is expecting JSON-encoded data on /search endpoint with these fields:
- title
- phrase

HTTPie example:
```shell
echo '{"title": "Romeo & Juliet", "phrase": "oh romeo romeo"}'  | http "http://localhost:8000/search" 
```

### Configuration

Application can be configured with environment variables, most important keys are presented below.
For more details and default values, please go to [config.go](src/cmd/api/config.go#L156).

```text
CACHE_ANSWER          # 0-1: enable/disable cache based on query sent to application and it's answer
CACHE_LISTING         # 0-1: enable/disable cache for listing metadata for given "title" part of query
CACHE_CONTENT         # 0-1: enable/disable cache for downloaded content (strongly suggested)
DOWNLOAD_DELAY_MIN    # download Task limitation to emulate human-like behaviour to prevent from banning,
DOWNLOAD_DELAY_MAX    # min/max value, each download gets random value from that range
SEARCH_WORKERS        # search worker goroutines (inefficient without cached content)
SEARCH_MAX_DISTANCE   # fuzzy-search engine distance option
SEARCH_RANDOM_RESULT  # returns a random match in the scope of given book instead of a first found match
                      #     [Note: cannot work properly with with CACHE_ANSWER enabled]
SEARCH_TIMEOUT        # maximum time allowed to spent by server for each search request
```

### tests

Bunch of tests are available, to run them simply run in `src/` directory:
```
go test ./... -v -cover
```

```
$ go test ./... -v -cover
  === RUN   Test_stringToTime
  === RUN   Test_stringToTime/input:'592ms'
  === RUN   Test_stringToTime/input:'595_ms'
  === RUN   Test_stringToTime/input:'12s'
  === RUN   Test_stringToTime/input:'15_s'
  === RUN   Test_stringToTime/input:'12m'
  === RUN   Test_stringToTime/input:'15_m'
  === RUN   Test_stringToTime/input:'12h'
  === RUN   Test_stringToTime/input:'15_h'
  --- PASS: Test_stringToTime (0.00s)
      --- PASS: Test_stringToTime/input:'592ms' (0.00s)
      --- PASS: Test_stringToTime/input:'595_ms' (0.00s)
      --- PASS: Test_stringToTime/input:'12s' (0.00s)
      --- PASS: Test_stringToTime/input:'15_s' (0.00s)
      --- PASS: Test_stringToTime/input:'12m' (0.00s)
      --- PASS: Test_stringToTime/input:'15_m' (0.00s)
      --- PASS: Test_stringToTime/input:'12h' (0.00s)
      --- PASS: Test_stringToTime/input:'15_h' (0.00s)
  === RUN   Test_stringToTimeErrors
  === RUN   Test_stringToTimeErrors/input:''
  === RUN   Test_stringToTimeErrors/input:'xd_s'
  === RUN   Test_stringToTimeErrors/input:'ms_254'
  --- PASS: Test_stringToTimeErrors (0.00s)
      --- PASS: Test_stringToTimeErrors/input:'' (0.00s)
      --- PASS: Test_stringToTimeErrors/input:'xd_s' (0.00s)
      --- PASS: Test_stringToTimeErrors/input:'ms_254' (0.00s)
  === RUN   Test_stringToBool
  === RUN   Test_stringToBool/input:'1'
  === RUN   Test_stringToBool/input:'true'
  === RUN   Test_stringToBool/input:'TrUe'
  === RUN   Test_stringToBool/input:'0'
  === RUN   Test_stringToBool/input:'false'
  === RUN   Test_stringToBool/input:'FaLsE'
  --- PASS: Test_stringToBool (0.00s)
      --- PASS: Test_stringToBool/input:'1' (0.00s)
      --- PASS: Test_stringToBool/input:'true' (0.00s)
      --- PASS: Test_stringToBool/input:'TrUe' (0.00s)
      --- PASS: Test_stringToBool/input:'0' (0.00s)
      --- PASS: Test_stringToBool/input:'false' (0.00s)
      --- PASS: Test_stringToBool/input:'FaLsE' (0.00s)
  === RUN   Test_stringToBoolErrors
  === RUN   Test_stringToBoolErrors/input:''
  === RUN   Test_stringToBoolErrors/input:'-1'
  === RUN   Test_stringToBoolErrors/input:'-0'
  === RUN   Test_stringToBoolErrors/input:'2'
  === RUN   Test_stringToBoolErrors/input:'999'
  === RUN   Test_stringToBoolErrors/input:'asdf'
  === RUN   Test_stringToBoolErrors/input:'QWERTY'
  --- PASS: Test_stringToBoolErrors (0.00s)
      --- PASS: Test_stringToBoolErrors/input:'' (0.00s)
      --- PASS: Test_stringToBoolErrors/input:'-1' (0.00s)
      --- PASS: Test_stringToBoolErrors/input:'-0' (0.00s)
      --- PASS: Test_stringToBoolErrors/input:'2' (0.00s)
      --- PASS: Test_stringToBoolErrors/input:'999' (0.00s)
      --- PASS: Test_stringToBoolErrors/input:'asdf' (0.00s)
      --- PASS: Test_stringToBoolErrors/input:'QWERTY' (0.00s)
  === RUN   Test_Search
  === RUN   Test_Search/0:Unexpected_format_of_payload_(not_JSON)
  2020/09/28 08:32:56 unmarshall failed: invalid character 'a' looking for beginning of value
  === RUN   Test_Search/1:Missing_fields
  === RUN   Test_Search/2:Missing_'phrase'_field
  === RUN   Test_Search/3:Missing_'title'_field
  === RUN   Test_Search/4:Payload_OK
  === RUN   Test_Search/5:Wrong_`title`_field_type
  2020/09/28 08:32:56 unmarshall failed: json: cannot unmarshal number into Go struct field Payload.title of type string
  === RUN   Test_Search/6:Wrong_`phrase`_field_type
  2020/09/28 08:32:56 unmarshall failed: json: cannot unmarshal bool into Go struct field Payload.phrase of type string
  --- PASS: Test_Search (0.00s)
      --- PASS: Test_Search/0:Unexpected_format_of_payload_(not_JSON) (0.00s)
      --- PASS: Test_Search/1:Missing_fields (0.00s)
      --- PASS: Test_Search/2:Missing_'phrase'_field (0.00s)
      --- PASS: Test_Search/3:Missing_'title'_field (0.00s)
      --- PASS: Test_Search/4:Payload_OK (0.00s)
      --- PASS: Test_Search/5:Wrong_`title`_field_type (0.00s)
      --- PASS: Test_Search/6:Wrong_`phrase`_field_type (0.00s)
  === RUN   Test_SearchError
  === RUN   Test_SearchError/0:Phrase_not_found
  === RUN   Test_SearchError/1:Processing_too_long
  --- PASS: Test_SearchError (0.00s)
      --- PASS: Test_SearchError/0:Phrase_not_found (0.00s)
      --- PASS: Test_SearchError/1:Processing_too_long (0.00s)
  PASS
  coverage: 45.9% of statements
  ok  	fuzzy-search/cmd/api	(cached)	coverage: 45.9% of statements
  ?   	fuzzy-search/internal/app/gutenbergsearch	[no test files]
  === RUN   TestExpectedContext
  --- PASS: TestExpectedContext (0.00s)
  PASS
  coverage: 77.8% of statements
  ok  	fuzzy-search/internal/pkg/context	(cached)	coverage: 77.8% of statements
  === RUN   TestFindBooks
  --- PASS: TestFindBooks (0.00s)
  === RUN   TestFindTxtLinkref
  --- PASS: TestFindTxtLinkref (0.00s)
  PASS
  coverage: 25.8% of statements
  ok  	fuzzy-search/internal/pkg/data	(cached)	coverage: 25.8% of statements
  ?   	fuzzy-search/internal/pkg/search	[no test files]
```

62.5% of total coverage at the moment, not great, not terrible ¯\\_(ツ)_/¯ 

### Current implementation concerns

There is one download task/goroutine at the time which receives download requests, the main reason behind this approach
is to emulate real user behaviour on required website, of course there is still some room to improvement but in general
there should not be problems with that.

However, given number of searching workers/goroutines (configuration) are spawned with every request, which is not a great design
for many concurrent request to this service. More work could be focused on this part.

Another thing that may be improved is fuzzy-searching approach itself. I never had experience with such mechanism 
before, and I'm more than sure it can be implemented in better ways. Solution that I wrote seems to work just fine
for need of this application.

More unit-tests would be also required in real software development scenario, as well as functional ones,
but for purpose of this task I did limit myself to bare minimum in that matter. I was more focused 
to have actually working application during its development.


### Example app logs

```
# answer cache enabled
CACHE_ANSWER=1 SEARCH_RANDOM_RESULT=1 ./build/local-run.sh 

2020/09/28 07:03:45 Loaded config:
2020/09/28 07:03:45 &main.Config{answerCache:true, answerCacheExpiration:14400000000000, answerCacheCleanupInterval:1860000000000, listingCache:true, listingCacheExpiration:14400000000000, listingCacheCleanupInterval:600000000000, contentCache:true, contentCacheExpiration:3600000000000, contentCacheCleanupInterval:600000000000, downloadDelayMin:1000000000, downloadDelayMax:2000000000, searchWorkers:8, searchMaxDistance:2, searchRandomResult:true, providerUserAgent:"Mozilla/5.0 (X11; Linux x86_64; rv:80.0) Gecko/20100101 Firefox/80.0", providerTimeout:0}
2020/09/28 07:03:45 Starting http server
2020/09/28 07:03:45 [[ downloadTask running ]]
2020/09/28 07:03:45 [[ searchTask running ]]

[stdin] echo '{"title": "Romeo & Juliet", "phrase": "o my romeo"}'  | http "http://localhost:8000/search"

2020/09/28 07:03:48 Searching books with "Romeo & Juliet" title
2020/09/28 07:03:50 Read 21 book positions from external source
2020/09/28 07:03:50 Scheduled 21 books to download
2020/09/28 07:03:50 [downloadTask] Download job acquired
2020/09/28 07:03:50 [searchTask] Search job acquired
2020/09/28 07:03:50 [searchTask] 8 search workers running
2020/09/28 07:03:50 [DWorker] sleeping for 1.555503445s
2020/09/28 07:03:53 [DWorker] Book ('The Tragedy of Romeo and Juliet' - William Shakespeare) downloaded in 3.239341891s
2020/09/28 07:03:53 [DWorker] sleeping for 1.978762601s
2020/09/28 07:03:53 result found! ('The Tragedy of Romeo and Juliet' - William Shakespeare)
2020/09/28 07:03:57 [DWorker] Book ('Romeo and Juliet' - William Shakespeare) downloaded in 7.205026184s
2020/09/28 07:03:57 [DWorker] Downloading interrupted
2020/09/28 07:03:57 [searchTask] search workers finishes

[stdout] HTTP/1.1 200 OK
[stdout] Content-Length: 201
[stdout] Content-Type: text/plain; charset=utf-8
[stdout] Date: Mon, 28 Sep 2020 05:03:53 GMT
[stdout] 
[stdout] O Romeo, Romeo!
[stdout]     Who ever would have thought it? Romeo!
[stdout] 
[stdout]   Jul. What devil art thou that dost torment me thus?
[stdout]     This torture should be roar'd in dismal hell.
[stdout]     Hath Romeo slain himself? Sa

[stdin] echo '{"title": "Romeo & Juliet", "phrase": "o my romeo"}'  | http "http://localhost:8000/search"
2020/09/28 07:04:59 found cached query result

[stdout] HTTP/1.1 200 OK
[stdout] Content-Length: 201
[stdout] Content-Type: text/plain; charset=utf-8
[stdout] Date: Mon, 28 Sep 2020 05:04:59 GMT
[stdout] 
[stdout] O Romeo, Romeo!
[stdout]     Who ever would have thought it? Romeo!
[stdout] 
[stdout]   Jul. What devil art thou that dost torment me thus?
[stdout]     This torture should be roar'd in dismal hell.
[stdout]     Hath Romeo slain himself? Sa
^C
interrupt
2020/09/28 07:06:36 Closing application...
2020/09/28 07:06:36 Application closed
2020/09/28 07:06:36 [[ downloadTask closed ]]
2020/09/28 07:06:36 Unexpected server error: http: Server closed
2020/09/28 07:06:36 Server stopped
2020/09/28 07:06:36 [[ searchTask closed ]]



# answer cache disabled
CACHE_ANSWER=0 SEARCH_RANDOM_RESULT=1 ./build/local-run.sh

2020/09/28 07:07:10 Loaded config:
2020/09/28 07:07:10 &main.Config{answerCache:false, answerCacheExpiration:14400000000000, answerCacheCleanupInterval:1860000000000, listingCache:true, listingCacheExpiration:14400000000000, listingCacheCleanupInterval:600000000000, contentCache:true, contentCacheExpiration:3600000000000, contentCacheCleanupInterval:600000000000, downloadDelayMin:1000000000, downloadDelayMax:2000000000, searchWorkers:8, searchMaxDistance:2, searchRandomResult:true, providerUserAgent:"Mozilla/5.0 (X11; Linux x86_64; rv:80.0) Gecko/20100101 Firefox/80.0", providerTimeout:0}
2020/09/28 07:07:10 Starting http server
2020/09/28 07:07:10 [[ downloadTask running ]]
2020/09/28 07:07:10 [[ searchTask running ]]

[stdin] echo '{"title": "Romeo & Juliet", "phrase": "o my romeo"}'  | http "http://localhost:8000/search"

2020/09/28 07:07:50 Searching books with "Romeo & Juliet" title
2020/09/28 07:07:50 Read 21 book positions from external source
2020/09/28 07:07:50 Scheduled 21 books to download
2020/09/28 07:07:50 [searchTask] Search job acquired
2020/09/28 07:07:50 [searchTask] 8 search workers running
2020/09/28 07:07:50 [downloadTask] Download job acquired
2020/09/28 07:07:50 [DWorker] sleeping for 1.732068003s
2020/09/28 07:07:54 [DWorker] Book ('The Tragedy of Romeo and Juliet' - William Shakespeare) downloaded in 4.076077373s
2020/09/28 07:07:54 [DWorker] sleeping for 1.075818762s
2020/09/28 07:07:55 result found! ('The Tragedy of Romeo and Juliet' - William Shakespeare)
2020/09/28 07:07:58 [DWorker] Book ('Romeo and Juliet' - William Shakespeare) downloaded in 7.589249776s
2020/09/28 07:07:58 [DWorker] Downloading interrupted
2020/09/28 07:07:58 [SWorker 3] Search interrupted
2020/09/28 07:07:58 [searchTask] search workers finishes

[stdout] HTTP/1.1 200 OK
[stdout] Content-Length: 201
[stdout] Content-Type: text/plain; charset=utf-8
[stdout] Date: Mon, 28 Sep 2020 05:07:55 GMT
[stdout] 
[stdout] O gentle Romeo,
[stdout]     If thou dost love, pronounce it faithfully.
[stdout]     Or if thou thinkest I am too quickly won,
[stdout]     I'll frown, and be perverse, and say thee nay,
[stdout]     So thou wilt woo; but else, not

[stdin] echo '{"title": "Romeo & Juliet", "phrase": "o my romeo"}'  | http "http://localhost:8000/search"

2020/09/28 07:08:22 Searching books with "Romeo & Juliet" title
2020/09/28 07:08:22 Read 21 book positions from cache
2020/09/28 07:08:22 Load book from cache ("The Tragedy of Romeo and Juliet" - William Shakespeare)
2020/09/28 07:08:22 Load book from cache ("Romeo and Juliet" - William Shakespeare)
2020/09/28 07:08:22 Scheduled 19 books to download
2020/09/28 07:08:22 [downloadTask] Download job acquired
2020/09/28 07:08:22 [searchTask] Search job acquired
2020/09/28 07:08:22 [searchTask] 8 search workers running
2020/09/28 07:08:22 [DWorker] sleeping for 1.04848212s
2020/09/28 07:08:22 result found! ('The Tragedy of Romeo and Juliet' - William Shakespeare)
2020/09/28 07:08:22 [SWorker 4] Search interrupted
2020/09/28 07:08:26 [DWorker] Book ('Beautiful Stories from Shakespeare' - William Shakespeare and E. Nesbit) downloaded in 3.731857445s
2020/09/28 07:08:26 [DWorker] Downloading interrupted
2020/09/28 07:08:26 [SWorker 3] Search interrupted
2020/09/28 07:08:26 [searchTask] search workers finishes

[stdout] HTTP/1.1 200 OK
[stdout] Content-Length: 201
[stdout] Content-Type: text/plain; charset=utf-8
[stdout] Date: Mon, 28 Sep 2020 05:08:22 GMT
[stdout] 
[stdout] O Romeo, Romeo, brave Mercutio's dead!
[stdout]     That gallant spirit hath aspir'd the clouds,
[stdout]     Which too untimely here did scorn the earth.
[stdout] 
[stdout]   Rom. This day's black fate on moe days doth depend;
```